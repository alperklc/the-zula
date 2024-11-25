package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alperklc/the-zula/service/infrastructure/environment"
	bookmarksService "github.com/alperklc/the-zula/service/services/bookmarks"
	importExportService "github.com/alperklc/the-zula/service/services/importExport"
	notesService "github.com/alperklc/the-zula/service/services/notes"
	userActivityService "github.com/alperklc/the-zula/service/services/userActivity"
	usersService "github.com/alperklc/the-zula/service/services/users"
	"github.com/gorilla/websocket"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
)

type a struct {
	env            *environment.Config
	users          usersService.UsersService
	userActivities userActivityService.UserActivityService
	bookmarks      bookmarksService.BookmarkService
	notes          notesService.NoteService
	importer       importExportService.ImportExportService
	clientHub      Hub
}

func NewApi(c *environment.Config, u usersService.UsersService, ua userActivityService.UserActivityService, bs bookmarksService.BookmarkService, n notesService.NoteService, is importExportService.ImportExportService, clientHub Hub) ServerInterface {
	return &a{env: c, users: u, userActivities: ua, bookmarks: bs, notes: n, importer: is, clientHub: clientHub}
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

func (s *a) GetFrontendConfig(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, FrontendConfig{
		Authority:             &s.env.FEAuthority,
		ClientId:              &s.env.FEClientId,
		PostLogoutRedirectUri: &s.env.FEPostLogoutRedirectUri,
		RedirectUri:           &s.env.FERedirectUri,
	})
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
	sessionId := r.Header.Get("sessionId")

	var input NoteInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create note, %s", err.Error()))
		return
	}

	noteCreated, errCreateNote := s.notes.CreateNote(user, sessionId, input.Title, input.Content, input.Tags)
	if errCreateNote != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not create note, %s", errCreateNote.Error()))
		return
	}

	sendResponse(w, http.StatusOK, noteCreated)
}

func (s *a) DeleteNote(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())
	sessionId := r.Header.Get("sessionId")

	err := s.notes.DeleteNote(id, user, sessionId)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not delete note, %s", err.Error()))
		return
	}

	sendResponse(w, http.StatusOK, "ok")
}

func (s *a) GetNote(w http.ResponseWriter, r *http.Request, id string, params GetNoteParams) {
	user := authorization.UserID(r.Context())
	sessionId := r.Header.Get("sessionId")

	p := notesService.GetNoteParams{}
	if params.LoadDraft != nil {
		p.LoadDraft = *params.LoadDraft
	}
	if params.GetChanges != nil {
		p.GetChanges = *params.GetChanges
	}
	if params.GetReferences != nil {
		p.GetReferences = *params.GetReferences
	}

	response, errGetNotes := s.notes.GetNote(id, user, sessionId, p)
	if errGetNotes != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not get note, %s", errGetNotes.Error()))
		return
	}

	sendResponse(w, http.StatusOK, response)
}

func (s *a) UpdateNote(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())
	sessionId := r.Header.Get("sessionId")

	var updateInput = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not update note, %s", err.Error()))
		return
	}

	errUpdateNote := s.notes.UpdateNote(id, user, sessionId, updateInput)
	if errUpdateNote != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not update note, %s", errUpdateNote.Error()))
		return
	}

	sendResponse(w, http.StatusOK, "ok")
}

func (s *a) GetNotesChanges(w http.ResponseWriter, r *http.Request, shortId string, params GetNotesChangesParams) {
	user := authorization.UserID(r.Context())

	response, errGetNotesChanges := s.notes.ListNotesChanges(user, shortId, params.Page, params.PageSize)
	if errGetNotesChanges != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not get notes changes, %s", errGetNotesChanges.Error()))
		return
	}

	converted := make([]NoteChange, 0, len(response.Items))
	for _, item := range response.Items {
		formattedUpdatedAt := item.UpdatedAt.Format(time.RFC3339)
		converted = append(converted, NoteChange{ShortId: &item.ShortId, NoteId: &item.NoteId, UpdatedAt: &formattedUpdatedAt, Change: &item.Change, UpdatedBy: &item.UpdatedBy})
	}

	sendResponse(w, http.StatusOK, NotesChangesResult{
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

func (s *a) GetNotesChange(w http.ResponseWriter, r *http.Request, noteId, shortId string) {
	response, errGetNotesChange := s.notes.GetNotesChange(noteId, shortId)
	if errGetNotesChange != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not get notes change, %s", errGetNotesChange.Error()))
		return
	}

	sendResponse(w, http.StatusOK, response)
}

func (s *a) GetBookmarks(w http.ResponseWriter, r *http.Request, params GetBookmarksParams) {
	user := authorization.UserID(r.Context())

	response, errGetBookmarks := s.bookmarks.ListBookmarks(user, params.Q, params.Page, params.PageSize, params.SortBy, params.SortDirection, params.Tags)
	if errGetBookmarks != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not get bookmarks, %s", errGetBookmarks.Error()))
		return
	}

	converted := make([]Bookmark, 0, len(response.Items))
	for _, item := range response.Items {
		converted = append(converted, Bookmark{ShortId: item.ShortId, UpdatedAt: item.UpdatedAt, Title: &item.Title, Tags: &item.Tags, Url: item.URL, FaviconUrl: &item.FaviconURL})
	}

	sendResponse(w, http.StatusOK, BookmarkSearchResult{
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

func (s *a) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	user := authorization.UserID(r.Context())
	sessionId := r.Header.Get("sessionId")

	var input BookmarkInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create bookmark, %s", err.Error()))
		return
	}

	noteCreated, errCreateNote := s.bookmarks.CreateBookmark(user, sessionId, input.Url, *input.Title, input.Tags)
	if errCreateNote != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not create bookmark, %s", errCreateNote.Error()))
		return
	}

	sendResponse(w, http.StatusOK, noteCreated)
}

func (s *a) DeleteBookmark(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())
	sessionId := r.Header.Get("sessionId")

	err := s.bookmarks.DeleteBookmark(id, user, sessionId)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not delete bookmark, %s", err.Error()))
		return
	}

	sendResponse(w, http.StatusOK, "ok")
}

func (s *a) GetBookmark(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())
	sessionId := r.Header.Get("sessionId")

	response, errGetNotes := s.bookmarks.GetBookmark(id, user, sessionId)
	if errGetNotes != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not get note, %s", errGetNotes.Error()))
		return
	}

	sendResponse(w, http.StatusOK, response)
}

func (s *a) UpdateBookmark(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())
	sessionId := r.Header.Get("sessionId")

	var updateInput = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not update bookmark, %s", err.Error()))
		return
	}

	errUpdateNote := s.bookmarks.UpdateBookmark(id, user, sessionId, updateInput)
	if errUpdateNote != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not update bookmark, %s", errUpdateNote.Error()))
		return
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

	if *params.Type == "note" {
		response, errGetTags := s.notes.SearchTags(user, searchKeyword, limit)
		if errGetTags != nil {
			sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not get tags, %s", errGetTags.Error()))
			return
		}

		sendResponse(w, http.StatusOK, response)
	} else {
		response, errGetTags := s.bookmarks.SearchTags(user, searchKeyword, limit)
		if errGetTags != nil {
			sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not get tags, %s", errGetTags.Error()))
			return
		}

		sendResponse(w, http.StatusOK, response)
	}
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
	user := authorization.UserID(r.Context())
	sessionId := r.Header.Get("sessionId")

	if user != id {
		sendErrorResponse(w, http.StatusForbidden, "could not get user: not allowed")
		return
	}

	var input UserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not create bookmark, %s", err.Error()))
		return
	}

	errUpdateUser := s.users.UpdateUser(id, sessionId, *input.Email, *input.Firstname, *input.Lastname, *input.Displayname, input.Language, input.Theme)
	if errUpdateUser != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not get user, %s", errUpdateUser.Error()))
		return
	}

	sendResponse(w, http.StatusOK, "ok")
}

func (s *a) GetUserActivity(w http.ResponseWriter, r *http.Request, id string, params GetUserActivityParams) {
	user := authorization.UserID(r.Context())

	response, errGetActivities := s.userActivities.List(user, *params.Page, *params.PageSize, *params.SortBy, *params.SortDirection)
	if errGetActivities != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not get users activities, %s", errGetActivities.Error()))
		return
	}

	converted := make([]UserActivity, 0, len(response.Items))
	for _, item := range response.Items {
		timestamp := item.Timestamp.UTC().Format(time.RFC3339)
		converted = append(converted, UserActivity{
			Action: &item.Action, ClientId: nil, ObjectId: &item.ObjectID, ResourceType: &item.ResourceType, Timestamp: &timestamp,
		})
	}

	sendResponse(w, http.StatusOK, UserActivityResult{
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

func (s *a) GetInsights(w http.ResponseWriter, r *http.Request, id string) {
	user := authorization.UserID(r.Context())

	activityGraph, mostVisited, lastVisited, nrOfNotes, nrOfBookmarks, errActivities := s.userActivities.GetInsightsForDashboard(user)
	if errActivities != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not get grouped activities, %s", errActivities.Error()))
		return
	}

	response := Insights{}
	response.ConvertInsights(activityGraph, mostVisited, lastVisited, nrOfNotes, nrOfBookmarks)
	sendResponse(w, http.StatusOK, response)
}

func (s *a) ConnectWs(w http.ResponseWriter, r *http.Request, user string) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Upgrading the HTTP connection socket connection
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	CreateNewSocketUser(&s.clientHub, connection, user)
}

func (s *a) ImportData(w http.ResponseWriter, r *http.Request) {
	file, _, errFormFile := r.FormFile("file")
	if errFormFile != nil {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("could not read file uploaded, %s", errFormFile))
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, errCopy := io.Copy(buf, file)
	if errCopy != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not copy file to buffer, %s", errCopy))
		return
	}

	result, importErr := s.importer.ProcessIncomingZipFile(buf.Bytes())
	if importErr != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("could not import data from the zip file, %s", importErr))
		return
	}

	sendResponse(w, http.StatusOK, result)
}

func (s *a) ExportData(w http.ResponseWriter, r *http.Request) {
	user := authorization.UserID(r.Context())

	file, errExport := s.importer.MakeZipFile(user)
	if errExport != nil {
		http.Error(w, fmt.Sprintf("could not create file, %s", errExport), http.StatusNotFound)
		return
	}

	file, err := os.Open("exports/zula-" + user + ".zip")
	if err != nil {
		http.Error(w, fmt.Sprintf("File not found %s", err), http.StatusNotFound)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to get file info %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name()))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Stream the file to the client
	_, err = io.Copy(w, file) // Sends the file content to the client
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send file %s", err), http.StatusInternalServerError)
		return
	}
}
