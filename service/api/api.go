package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	notectrl "github.com/alperklc/the-zula/service/controller/notes"
	notesReferencesCtrl "github.com/alperklc/the-zula/service/controller/notesReferences"
)

type a struct {
	notes           notectrl.NoteController
	notesReferences notesReferencesCtrl.NotesReferencesController
}

func NewApi(n notectrl.NoteController, nr notesReferencesCtrl.NotesReferencesController) ServerInterface {
	return &a{notes: n, notesReferences: nr}
}

func (s *a) GetApiV1Me(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *a) DeleteApiV1NoteIdDraft(w http.ResponseWriter, r *http.Request, id string) {
	user := r.Context().Value("user").(string)
	w.Header().Set("Content-Type", "application/json")

	err := s.notes.DeleteDraft(user, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("could not delete draft, %s", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("ok")
}

func (s *a) PutApiV1NoteIdDraft(w http.ResponseWriter, r *http.Request, id string) {
	user := r.Context().Value("user").(string)
	w.Header().Set("Content-Type", "application/json")

	var requestBody NoteInput
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || *requestBody.Title == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("could not create draft, %s", err.Error()))
		return
	}

	errUpdate := s.notes.UpdateDraft(user, id, *requestBody.Title, *requestBody.Content, *requestBody.Tags)
	if errUpdate != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("could not create draft, %s", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("ok")
}

func (s *a) GetApiV1Notes(w http.ResponseWriter, r *http.Request, params GetApiV1NotesParams) {
	user := r.Context().Value("user").(string)
	w.Header().Set("Content-Type", "application/json")

	response, errGetNotes := s.notes.ListNotes(user, params.Q, params.Page, params.PageSize, params.SortBy, params.SortDirection, params.Tags)
	if errGetNotes != nil {
		http.Error(w, errGetNotes.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("could not get notes, %s", errGetNotes.Error()))
		return
	}

	converted := make([]Note, 0, len(response.Items))
	for _, item := range response.Items {
		formattedUpdatedAt := item.UpdatedAt.Format(time.RFC3339)
		converted = append(converted, Note{Id: &item.Id, UpdatedAt: &formattedUpdatedAt, Title: &item.Title, Tags: &item.Tags, HasDraft: &item.HasDraft})
	}

	json.NewEncoder(w).Encode(NoteSearchResult{
		Meta: &PaginationMeta{
			Count:         &response.Meta.Count,
			Page:          &response.Meta.Page,
			PageSize:      &response.Meta.PageSize,
			SortBy:        &response.Meta.SortBy,
			SortDirection: &response.Meta.SortDirection,
		},
		Items: &converted,
	})
}

func (s *a) PostApiV1Notes(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)
	w.Header().Set("Content-Type", "application/json")

	var input NoteInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("could not create note, %s", err.Error()))
		return
	}

	noteCreated, errCreateNote := s.notes.CreateNote(user, "1", input.Title, input.Content, input.Tags)
	if errCreateNote != nil {
		http.Error(w, errCreateNote.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(noteCreated)
}

func (s *a) DeleteApiV1NotesId(w http.ResponseWriter, r *http.Request, id string) {
	user := r.Context().Value("user").(string)
	w.Header().Set("Content-Type", "application/json")

	err := s.notes.DeleteNote(id, user, "1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("could not delete note, %s", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("ok")
}

func (s *a) GetApiV1NotesId(w http.ResponseWriter, r *http.Request, id string, params GetApiV1NotesIdParams) {
	user := r.Context().Value("user").(string)
	w.Header().Set("Content-Type", "application/json")

	response, errGetNotes := s.notes.GetNote(id, user, "1")
	if errGetNotes != nil {
		http.Error(w, errGetNotes.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("could not get notes, %s", errGetNotes.Error()))
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (s *a) PutApiV1NotesId(w http.ResponseWriter, r *http.Request, id string) {
	user := r.Context().Value("user").(string)
	w.Header().Set("Content-Type", "application/json")

	var updateInput = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("could not create note, %s", err.Error()))
		return
	}

	errUpdateNote := s.notes.UpdateNote(id, user, "1", updateInput)
	if errUpdateNote != nil {
		http.Error(w, errUpdateNote.Error(), http.StatusInternalServerError)
		return
	}

	content, contentChanged := updateInput["content"]
	if contentChanged {
		s.notesReferences.UpsertReferencesOfNote(id, content.(string))
	}

	json.NewEncoder(w).Encode("ok")
}

func (s *a) GetApiV1Tags(w http.ResponseWriter, r *http.Request, params GetApiV1TagsParams) {
	user := r.Context().Value("user").(string)
	w.Header().Set("Content-Type", "application/json")

	var limit int = 10
	l := r.URL.Query().Get("limit")
	if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 && parsedLimit < 10 {
		limit = parsedLimit
	}

	var searchKeyword string = ""
	if q := r.URL.Query().Get("q"); q != "" {
		searchKeyword = q
	}

	response, errGetTags := s.notes.SearchTags(user, searchKeyword, limit)
	if errGetTags != nil {
		http.Error(w, errGetTags.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("could not get tags, %s", errGetTags.Error()))
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (s *a) GetApiV1UsersId(w http.ResponseWriter, r *http.Request, id string) {

}

func (s *a) GetApiV1UsersIdActivity(w http.ResponseWriter, r *http.Request, id string, params GetApiV1UsersIdActivityParams) {

}

func (s *a) GetApiV1UsersIdInsights(w http.ResponseWriter, r *http.Request, id string) {

}
