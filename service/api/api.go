package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	notectrl "github.com/alperklc/the-zula/service/services/notes"
	notesReferencesCtrl "github.com/alperklc/the-zula/service/services/notesReferences"
	"github.com/alperklc/the-zula/service/services/users"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
)

type a struct {
	users           users.UsersService
	notes           notectrl.NoteService
	notesReferences notesReferencesCtrl.NotesReferencesService
}

func NewApi(u users.UsersService, n notectrl.NoteService, nr notesReferencesCtrl.NotesReferencesService) ServerInterface {
	return &a{users: u, notes: n, notesReferences: nr}
}

func sendResponse(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func sendErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(APIErrorResponse{Message: message})
}

func (s *a) DeleteNoteDraft(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())

	err := s.notes.DeleteDraft(user, id)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not delete draft, %s", err.Error()))
		return
	}

	sendResponse(w, http.StatusOK, "ok")
}

func (s *a) SaveNoteDraft(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())

	var requestBody NoteInput
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || *requestBody.Title == "" {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create draft, %s", err.Error()))
		return
	}

	errUpdate := s.notes.UpdateDraft(user, id, *requestBody.Title, *requestBody.Content, *requestBody.Tags)
	if errUpdate != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create draft, %s", errUpdate.Error()))
		return
	}

	sendResponse(w, http.StatusOK, "ok")
}

func (s *a) GetNotes(w http.ResponseWriter, r *http.Request, params GetNotesParams) {
	user := authorization.UserID(r.Context())

	response, errGetNotes := s.notes.ListNotes(user, params.Q, params.Page, params.PageSize, params.SortBy, params.SortDirection, params.Tags)
	if errGetNotes != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not get notes, %s", errGetNotes.Error()))
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

func (s *a) CreateNote(w http.ResponseWriter, r *http.Request) {
	user := authorization.UserID(r.Context())

	var input NoteInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create note, %s", err.Error()))
		return
	}

	noteCreated, errCreateNote := s.notes.CreateNote(user, "1", input.Title, input.Content, input.Tags)
	if errCreateNote != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not create note, %s", errCreateNote.Error()))
		return
	}

	go s.notesReferences.UpsertReferencesOfNote(noteCreated.ShortId, *input.Content)

	sendResponse(w, http.StatusOK, noteCreated)
}

func (s *a) DeleteNote(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())

	err := s.notes.DeleteNote(id, user, "1")
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not delete note, %s", err.Error()))
		return
	}

	sendResponse(w, http.StatusOK, "ok")
}

func (s *a) GetNote(w http.ResponseWriter, r *http.Request, id string, params GetNoteParams) {
	user := authorization.UserID(r.Context())

	response, errGetNotes := s.notes.GetNote(id, user, "1")
	if errGetNotes != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not get note, %s", errGetNotes.Error()))
		return
	}

	if params.LoadDraft != nil && *params.LoadDraft {
		draftOfNote, errGetDraft := s.notes.GetDraftOfNote(user, id)
		if errGetDraft != nil {
			sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not get draft of the note, %s", errGetDraft.Error()))
			return
		}

		response.Content = draftOfNote.Content
		response.Title = draftOfNote.Title
		response.Tags = draftOfNote.Tags
	}

	sendResponse(w, http.StatusOK, response)
}

func (s *a) UpdateNote(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())

	var updateInput = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create note, %s", err.Error()))
		return
	}

	errUpdateNote := s.notes.UpdateNote(id, user, "1", updateInput)
	if errUpdateNote != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create note, %s", errUpdateNote.Error()))
		return
	}

	content, contentChanged := updateInput["content"]
	if contentChanged {
		go s.notesReferences.UpsertReferencesOfNote(id, content.(string))
	}

	sendResponse(w, http.StatusOK, "ok")
}

func (s *a) GetTags(w http.ResponseWriter, r *http.Request, params GetTagsParams) {
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
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not get tags, %s", errGetTags.Error()))
		return
	}

	sendResponse(w, http.StatusOK, response)
}

func (s *a) GetUser(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())
	if user != id {
		sendErrorResponse(w, http.StatusForbidden, "could not get user: not allowed")
		return
	}

	response, errGetUser := s.users.GetUser(id)
	if errGetUser != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not get user, %s", errGetUser.Error()))
		return
	}

	sendResponse(w, http.StatusOK, response)
}

func (s *a) UpdateUser(w http.ResponseWriter, r *http.Request, id string) {

}

func (s *a) GetUserActivity(w http.ResponseWriter, r *http.Request, id string, params GetUserActivityParams) {

}

func (s *a) GetInsights(w http.ResponseWriter, r *http.Request, id string) {

}
