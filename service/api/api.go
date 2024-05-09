package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	notectrl "github.com/alperklc/the-zula/service/controller/notes"
	notesReferencesCtrl "github.com/alperklc/the-zula/service/controller/notesReferences"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
)

type a struct {
	notes           notectrl.NoteController
	notesReferences notesReferencesCtrl.NotesReferencesController
}

func NewApi(n notectrl.NoteController, nr notesReferencesCtrl.NotesReferencesController) ServerInterface {
	return &a{notes: n, notesReferences: nr}
}

func sendResponse(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func (s *a) GetApiV1Me(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *a) DeleteApiV1NoteShortIdDraft(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())

	err := s.notes.DeleteDraft(user, id)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, fmt.Sprintf("could not delete draft, %s", err.Error()))
		return
	}

	sendResponse(w, http.StatusOK, "ok")
}

func (s *a) PutApiV1NoteShortIdDraft(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())

	var requestBody NoteInput
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || *requestBody.Title == "" {
		sendResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create draft, %s", err.Error()))
		return
	}

	errUpdate := s.notes.UpdateDraft(user, id, *requestBody.Title, *requestBody.Content, *requestBody.Tags)
	if errUpdate != nil {
		sendResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create draft, %s", errUpdate.Error()))
		return
	}

	sendResponse(w, http.StatusOK, "ok")
}

func (s *a) GetApiV1Notes(w http.ResponseWriter, r *http.Request, params GetApiV1NotesParams) {
	user := authorization.UserID(r.Context())

	response, errGetNotes := s.notes.ListNotes(user, params.Q, params.Page, params.PageSize, params.SortBy, params.SortDirection, params.Tags)
	if errGetNotes != nil {
		sendResponse(w, http.StatusBadRequest, fmt.Sprintf("could not get notes, %s", errGetNotes.Error()))
		return
	}

	converted := make([]Note, 0, len(response.Items))
	for _, item := range response.Items {
		formattedUpdatedAt := item.UpdatedAt.Format(time.RFC3339)
		converted = append(converted, Note{ShortId: &item.ShortId, UpdatedAt: &formattedUpdatedAt, Title: &item.Title, Tags: &item.Tags, HasDraft: &item.HasDraft})
	}

	sendResponse(w, http.StatusOK, NoteSearchResult{
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
	user := authorization.UserID(r.Context())

	var input NoteInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create note, %s", err.Error()))
		return
	}

	noteCreated, errCreateNote := s.notes.CreateNote(user, "1", input.Title, input.Content, input.Tags)
	if errCreateNote != nil {
		sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not create note, %s", errCreateNote.Error()))
		return
	}

	go s.notesReferences.UpsertReferencesOfNote(noteCreated.ShortId, *input.Content)

	sendResponse(w, http.StatusOK, noteCreated)
}

func (s *a) DeleteApiV1NotesShortId(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())

	err := s.notes.DeleteNote(id, user, "1")
	if err != nil {
		sendResponse(w, http.StatusBadRequest, fmt.Sprintf("could not delete note, %s", err.Error()))
		return
	}

	sendResponse(w, http.StatusOK, "ok")
}

func (s *a) GetApiV1NotesShortId(w http.ResponseWriter, r *http.Request, id string, params GetApiV1NotesShortIdParams) {
	user := authorization.UserID(r.Context())

	response, errGetNotes := s.notes.GetNote(id, user, "1")
	if errGetNotes != nil {
		sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not get notes, %s", errGetNotes.Error()))
		return
	}

	if params.LoadDraft != nil && *params.LoadDraft {
		draftOfNote, errGetDraft := s.notes.GetDraftOfNote(user, id)
		if errGetDraft != nil {
			sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not get draft of the note, %s", errGetDraft.Error()))
			return
		}

		response.Content = draftOfNote.Content
		response.Title = draftOfNote.Title
		response.Tags = draftOfNote.Tags
	}

	sendResponse(w, http.StatusOK, response)
}

func (s *a) PutApiV1NotesShortId(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())

	var updateInput = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create note, %s", err.Error()))
		return
	}

	errUpdateNote := s.notes.UpdateNote(id, user, "1", updateInput)
	if errUpdateNote != nil {
		sendResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create note, %s", errUpdateNote.Error()))
		return
	}

	content, contentChanged := updateInput["content"]
	if contentChanged {
		go s.notesReferences.UpsertReferencesOfNote(id, content.(string))
	}

	sendResponse(w, http.StatusOK, "ok")
}

func (s *a) GetApiV1Tags(w http.ResponseWriter, r *http.Request, params GetApiV1TagsParams) {
	user := authorization.UserID(r.Context())

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
		sendResponse(w, http.StatusBadRequest, fmt.Sprintf("could not get tags, %s", errGetTags.Error()))
		return
	}

	sendResponse(w, http.StatusOK, response)
}

func (s *a) GetApiV1UsersShortId(w http.ResponseWriter, r *http.Request, id string) {

}

func (s *a) GetApiV1UsersShortIdActivity(w http.ResponseWriter, r *http.Request, id string, params GetApiV1UsersShortIdActivityParams) {

}

func (s *a) GetApiV1UsersShortIdInsights(w http.ResponseWriter, r *http.Request, id string) {

}
