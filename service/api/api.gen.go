// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ActivityOnDate defines model for activityOnDate.
type ActivityOnDate struct {
	Count *float32            `json:"count,omitempty"`
	Date  *openapi_types.Date `json:"date,omitempty"`
}

// Bookmark defines model for bookmark.
type Bookmark struct {
	CreatedAt   time.Time    `json:"createdAt"`
	CreatedBy   *string      `json:"createdBy,omitempty"`
	FaviconUrl  *string      `json:"faviconUrl,omitempty"`
	PageContent *PageContent `json:"pageContent,omitempty"`
	ShortId     string       `json:"shortId"`
	Tags        *[]string    `json:"tags,omitempty"`
	Title       *string      `json:"title,omitempty"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	UpdatedBy   *string      `json:"updatedBy,omitempty"`
	Url         string       `json:"url"`
}

// BookmarkInput defines model for bookmarkInput.
type BookmarkInput struct {
	Tags  *[]string `json:"tags,omitempty"`
	Title *string   `json:"title,omitempty"`
	Url   string    `json:"url"`
}

// BookmarkSearchResult defines model for bookmarkSearchResult.
type BookmarkSearchResult struct {
	Items *[]Bookmark     `json:"items,omitempty"`
	Meta  *PaginationMeta `json:"meta,omitempty"`
}

// Insights defines model for insights.
type Insights struct {
	ActivityGraph     *[]ActivityOnDate     `json:"activityGraph,omitempty"`
	LastVisited       *[]VisitingStatistics `json:"lastVisited,omitempty"`
	MostVisited       *[]MostVisited        `json:"mostVisited,omitempty"`
	NumberOfBookmarks *float32              `json:"numberOfBookmarks,omitempty"`
	NumberOfNotes     *float32              `json:"numberOfNotes,omitempty"`
}

// MostVisited defines model for mostVisited.
type MostVisited struct {
	Count    *float32 `json:"count,omitempty"`
	Id       *string  `json:"id,omitempty"`
	Name     *string  `json:"name,omitempty"`
	Title    *string  `json:"title,omitempty"`
	Typename *string  `json:"typename,omitempty"`
}

// Note defines model for note.
type Note struct {
	Content    *string         `json:"content,omitempty"`
	CreatedAt  *string         `json:"createdAt,omitempty"`
	CreatedBy  *string         `json:"createdBy,omitempty"`
	HasDraft   *bool           `json:"hasDraft,omitempty"`
	References *NoteReferences `json:"references,omitempty"`
	ShortId    *string         `json:"shortId,omitempty"`
	Tags       *[]string       `json:"tags,omitempty"`
	Title      *string         `json:"title,omitempty"`
	UpdatedAt  *string         `json:"updatedAt,omitempty"`
	UpdatedBy  *string         `json:"updatedBy,omitempty"`
}

// NoteInput defines model for noteInput.
type NoteInput struct {
	Content *string   `json:"content,omitempty"`
	Tags    *[]string `json:"tags,omitempty"`
	Title   *string   `json:"title,omitempty"`
}

// NoteLite defines model for noteLite.
type NoteLite struct {
	ShortId *string `json:"shortId,omitempty"`
	Title   *string `json:"title,omitempty"`
}

// NoteReferenceLink defines model for noteReferenceLink.
type NoteReferenceLink struct {
	Source *string `json:"source,omitempty"`
	Target *string `json:"target,omitempty"`
}

// NoteReferences defines model for noteReferences.
type NoteReferences struct {
	Links *[]NoteReferenceLink `json:"links,omitempty"`
	Meta  *PaginationMeta      `json:"meta,omitempty"`
	Nodes *[]NoteLite          `json:"nodes,omitempty"`
}

// NoteSearchResult defines model for noteSearchResult.
type NoteSearchResult struct {
	Items *[]Note         `json:"items,omitempty"`
	Meta  *PaginationMeta `json:"meta,omitempty"`
}

// PageContent defines model for pageContent.
type PageContent struct {
	Author    *string `json:"author,omitempty"`
	Favicon   *string `json:"favicon,omitempty"`
	Image     *string `json:"image,omitempty"`
	Length    *int    `json:"length,omitempty"`
	MdContent *string `json:"mdContent,omitempty"`
	SiteName  *string `json:"siteName,omitempty"`
	Title     *string `json:"title,omitempty"`
	Url       *string `json:"url,omitempty"`
}

// PaginationMeta defines model for paginationMeta.
type PaginationMeta struct {
	Count         *int    `json:"count,omitempty"`
	HasNextPage   *bool   `json:"hasNextPage,omitempty"`
	Page          *int    `json:"page,omitempty"`
	PageSize      *int    `json:"pageSize,omitempty"`
	Range         *string `json:"range,omitempty"`
	SortBy        *string `json:"sortBy,omitempty"`
	SortDirection *string `json:"sortDirection,omitempty"`
}

// Tag defines model for tag.
type Tag struct {
	Frequency    *int    `json:"frequency,omitempty"`
	TypeOfParent *string `json:"typeOfParent,omitempty"`
	Value        *string `json:"value,omitempty"`
}

// User defines model for user.
type User struct {
	CreatedAt   *string `json:"createdAt,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	Email       *string `json:"email,omitempty"`
	FirstName   *string `json:"firstName,omitempty"`
	Language    *string `json:"language,omitempty"`
	LastName    *string `json:"lastName,omitempty"`
	ShortId     *string `json:"shortId,omitempty"`
	Theme       *string `json:"theme,omitempty"`
}

// UserActivity defines model for userActivity.
type UserActivity struct {
	Action       *string `json:"action,omitempty"`
	ClientId     *string `json:"clientId,omitempty"`
	ObjectId     *string `json:"objectId,omitempty"`
	ResourceType *string `json:"resourceType,omitempty"`
	Timestamp    *string `json:"timestamp,omitempty"`
}

// UserActivityResult defines model for userActivityResult.
type UserActivityResult struct {
	Items *[]UserActivity `json:"items,omitempty"`
	Meta  *PaginationMeta `json:"meta,omitempty"`
}

// UserInput defines model for userInput.
type UserInput struct {
	CreatedAt *string `json:"createdAt,omitempty"`
	Email     *string `json:"email,omitempty"`
	Fullname  *string `json:"fullname,omitempty"`
	Language  *string `json:"language,omitempty"`
	Theme     *string `json:"theme,omitempty"`
	Username  *string `json:"username,omitempty"`
}

// VisitingStatistics defines model for visitingStatistics.
type VisitingStatistics struct {
	Id       *string `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	Title    *string `json:"title,omitempty"`
	Typename *string `json:"typename,omitempty"`
}

// GetBookmarksParams defines parameters for GetBookmarks.
type GetBookmarksParams struct {
	Q             *string   `form:"q,omitempty" json:"q,omitempty"`
	Page          *int      `form:"page,omitempty" json:"page,omitempty"`
	PageSize      *int      `form:"pageSize,omitempty" json:"pageSize,omitempty"`
	SortBy        *string   `form:"sortBy,omitempty" json:"sortBy,omitempty"`
	SortDirection *string   `form:"sortDirection,omitempty" json:"sortDirection,omitempty"`
	Tags          *[]string `form:"tags,omitempty" json:"tags,omitempty"`
}

// GetNotesParams defines parameters for GetNotes.
type GetNotesParams struct {
	Q             *string   `form:"q,omitempty" json:"q,omitempty"`
	Page          *int      `form:"page,omitempty" json:"page,omitempty"`
	PageSize      *int      `form:"pageSize,omitempty" json:"pageSize,omitempty"`
	SortBy        *string   `form:"sortBy,omitempty" json:"sortBy,omitempty"`
	SortDirection *string   `form:"sortDirection,omitempty" json:"sortDirection,omitempty"`
	Tags          *[]string `form:"tags,omitempty" json:"tags,omitempty"`
}

// GetNoteParams defines parameters for GetNote.
type GetNoteParams struct {
	LoadDraft *bool `form:"loadDraft,omitempty" json:"loadDraft,omitempty"`
}

// GetTagsParams defines parameters for GetTags.
type GetTagsParams struct {
	Type *string `form:"type,omitempty" json:"type,omitempty"`
	Q    *string `form:"q,omitempty" json:"q,omitempty"`
}

// GetUserActivityParams defines parameters for GetUserActivity.
type GetUserActivityParams struct {
	Page          *int    `form:"page,omitempty" json:"page,omitempty"`
	PageSize      *int    `form:"pageSize,omitempty" json:"pageSize,omitempty"`
	SortBy        *string `form:"sortBy,omitempty" json:"sortBy,omitempty"`
	SortDirection *string `form:"sortDirection,omitempty" json:"sortDirection,omitempty"`
}

// CreateBookmarkJSONRequestBody defines body for CreateBookmark for application/json ContentType.
type CreateBookmarkJSONRequestBody = BookmarkInput

// UpdateBookmarkJSONRequestBody defines body for UpdateBookmark for application/json ContentType.
type UpdateBookmarkJSONRequestBody = BookmarkInput

// CreateNoteJSONRequestBody defines body for CreateNote for application/json ContentType.
type CreateNoteJSONRequestBody = NoteInput

// UpdateNoteJSONRequestBody defines body for UpdateNote for application/json ContentType.
type UpdateNoteJSONRequestBody = NoteInput

// SaveNoteDraftJSONRequestBody defines body for SaveNoteDraft for application/json ContentType.
type SaveNoteDraftJSONRequestBody = NoteInput

// UpdateUserJSONRequestBody defines body for UpdateUser for application/json ContentType.
type UpdateUserJSONRequestBody = UserInput

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List bookmarks
	// (GET /api/v1/bookmarks)
	GetBookmarks(w http.ResponseWriter, r *http.Request, params GetBookmarksParams)
	// Create a new bookmark
	// (POST /api/v1/bookmarks)
	CreateBookmark(w http.ResponseWriter, r *http.Request)
	// Delete a bookmark by shortId
	// (DELETE /api/v1/bookmarks/{shortId})
	DeleteBookmark(w http.ResponseWriter, r *http.Request, shortId string)
	// Get a bookmark by shortId
	// (GET /api/v1/bookmarks/{shortId})
	GetBookmark(w http.ResponseWriter, r *http.Request, shortId string)
	// Update a bookmark by shortId
	// (PUT /api/v1/bookmarks/{shortId})
	UpdateBookmark(w http.ResponseWriter, r *http.Request, shortId string)
	// List notes
	// (GET /api/v1/notes)
	GetNotes(w http.ResponseWriter, r *http.Request, params GetNotesParams)
	// Create a new note
	// (POST /api/v1/notes)
	CreateNote(w http.ResponseWriter, r *http.Request)
	// Delete a note by shortId
	// (DELETE /api/v1/notes/{shortId})
	DeleteNote(w http.ResponseWriter, r *http.Request, shortId string)
	// Get a note by shortId
	// (GET /api/v1/notes/{shortId})
	GetNote(w http.ResponseWriter, r *http.Request, shortId string, params GetNoteParams)
	// Update a note by shortId
	// (PUT /api/v1/notes/{shortId})
	UpdateNote(w http.ResponseWriter, r *http.Request, shortId string)
	// Delete a notes draft by note shortId
	// (DELETE /api/v1/notes/{shortId}/draft)
	DeleteNoteDraft(w http.ResponseWriter, r *http.Request, shortId string)
	// Save draft of a note by notes shortId
	// (PUT /api/v1/notes/{shortId}/draft)
	SaveNoteDraft(w http.ResponseWriter, r *http.Request, shortId string)
	// Get tags
	// (GET /api/v1/tags)
	GetTags(w http.ResponseWriter, r *http.Request, params GetTagsParams)
	// Get a user by ID
	// (GET /api/v1/users/{shortId})
	GetUser(w http.ResponseWriter, r *http.Request, shortId string)
	// Update a user by ID
	// (PUT /api/v1/users/{shortId})
	UpdateUser(w http.ResponseWriter, r *http.Request, shortId string)
	// Get user activity
	// (GET /api/v1/users/{shortId}/activity)
	GetUserActivity(w http.ResponseWriter, r *http.Request, shortId string, params GetUserActivityParams)
	// Get dashboard insights
	// (GET /api/v1/users/{shortId}/insights)
	GetInsights(w http.ResponseWriter, r *http.Request, shortId string)
	// Connect to websocket
	// (GET /api/v1/ws/{user})
	ConnectWs(w http.ResponseWriter, r *http.Request, user string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetBookmarks operation middleware
func (siw *ServerInterfaceWrapper) GetBookmarks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetBookmarksParams

	// ------------- Optional query parameter "q" -------------

	err = runtime.BindQueryParameter("form", true, false, "q", r.URL.Query(), &params.Q)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "q", Err: err})
		return
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page", Err: err})
		return
	}

	// ------------- Optional query parameter "pageSize" -------------

	err = runtime.BindQueryParameter("form", true, false, "pageSize", r.URL.Query(), &params.PageSize)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "pageSize", Err: err})
		return
	}

	// ------------- Optional query parameter "sortBy" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortBy", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sortBy", Err: err})
		return
	}

	// ------------- Optional query parameter "sortDirection" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortDirection", r.URL.Query(), &params.SortDirection)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sortDirection", Err: err})
		return
	}

	// ------------- Optional query parameter "tags" -------------

	err = runtime.BindQueryParameter("form", true, false, "tags", r.URL.Query(), &params.Tags)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "tags", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetBookmarks(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateBookmark operation middleware
func (siw *ServerInterfaceWrapper) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateBookmark(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// DeleteBookmark operation middleware
func (siw *ServerInterfaceWrapper) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortId" -------------
	var shortId string

	err = runtime.BindStyledParameterWithOptions("simple", "shortId", mux.Vars(r)["shortId"], &shortId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteBookmark(w, r, shortId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetBookmark operation middleware
func (siw *ServerInterfaceWrapper) GetBookmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortId" -------------
	var shortId string

	err = runtime.BindStyledParameterWithOptions("simple", "shortId", mux.Vars(r)["shortId"], &shortId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetBookmark(w, r, shortId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// UpdateBookmark operation middleware
func (siw *ServerInterfaceWrapper) UpdateBookmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortId" -------------
	var shortId string

	err = runtime.BindStyledParameterWithOptions("simple", "shortId", mux.Vars(r)["shortId"], &shortId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateBookmark(w, r, shortId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetNotes operation middleware
func (siw *ServerInterfaceWrapper) GetNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetNotesParams

	// ------------- Optional query parameter "q" -------------

	err = runtime.BindQueryParameter("form", true, false, "q", r.URL.Query(), &params.Q)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "q", Err: err})
		return
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page", Err: err})
		return
	}

	// ------------- Optional query parameter "pageSize" -------------

	err = runtime.BindQueryParameter("form", true, false, "pageSize", r.URL.Query(), &params.PageSize)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "pageSize", Err: err})
		return
	}

	// ------------- Optional query parameter "sortBy" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortBy", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sortBy", Err: err})
		return
	}

	// ------------- Optional query parameter "sortDirection" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortDirection", r.URL.Query(), &params.SortDirection)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sortDirection", Err: err})
		return
	}

	// ------------- Optional query parameter "tags" -------------

	err = runtime.BindQueryParameter("form", true, false, "tags", r.URL.Query(), &params.Tags)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "tags", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetNotes(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateNote operation middleware
func (siw *ServerInterfaceWrapper) CreateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateNote(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// DeleteNote operation middleware
func (siw *ServerInterfaceWrapper) DeleteNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortId" -------------
	var shortId string

	err = runtime.BindStyledParameterWithOptions("simple", "shortId", mux.Vars(r)["shortId"], &shortId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteNote(w, r, shortId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetNote operation middleware
func (siw *ServerInterfaceWrapper) GetNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortId" -------------
	var shortId string

	err = runtime.BindStyledParameterWithOptions("simple", "shortId", mux.Vars(r)["shortId"], &shortId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortId", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetNoteParams

	// ------------- Optional query parameter "loadDraft" -------------

	err = runtime.BindQueryParameter("form", true, false, "loadDraft", r.URL.Query(), &params.LoadDraft)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "loadDraft", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetNote(w, r, shortId, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// UpdateNote operation middleware
func (siw *ServerInterfaceWrapper) UpdateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortId" -------------
	var shortId string

	err = runtime.BindStyledParameterWithOptions("simple", "shortId", mux.Vars(r)["shortId"], &shortId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateNote(w, r, shortId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// DeleteNoteDraft operation middleware
func (siw *ServerInterfaceWrapper) DeleteNoteDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortId" -------------
	var shortId string

	err = runtime.BindStyledParameterWithOptions("simple", "shortId", mux.Vars(r)["shortId"], &shortId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteNoteDraft(w, r, shortId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// SaveNoteDraft operation middleware
func (siw *ServerInterfaceWrapper) SaveNoteDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortId" -------------
	var shortId string

	err = runtime.BindStyledParameterWithOptions("simple", "shortId", mux.Vars(r)["shortId"], &shortId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SaveNoteDraft(w, r, shortId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetTags operation middleware
func (siw *ServerInterfaceWrapper) GetTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetTagsParams

	// ------------- Optional query parameter "type" -------------

	err = runtime.BindQueryParameter("form", true, false, "type", r.URL.Query(), &params.Type)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "type", Err: err})
		return
	}

	// ------------- Optional query parameter "q" -------------

	err = runtime.BindQueryParameter("form", true, false, "q", r.URL.Query(), &params.Q)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "q", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTags(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetUser operation middleware
func (siw *ServerInterfaceWrapper) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortId" -------------
	var shortId string

	err = runtime.BindStyledParameterWithOptions("simple", "shortId", mux.Vars(r)["shortId"], &shortId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetUser(w, r, shortId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// UpdateUser operation middleware
func (siw *ServerInterfaceWrapper) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortId" -------------
	var shortId string

	err = runtime.BindStyledParameterWithOptions("simple", "shortId", mux.Vars(r)["shortId"], &shortId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateUser(w, r, shortId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetUserActivity operation middleware
func (siw *ServerInterfaceWrapper) GetUserActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortId" -------------
	var shortId string

	err = runtime.BindStyledParameterWithOptions("simple", "shortId", mux.Vars(r)["shortId"], &shortId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortId", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetUserActivityParams

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page", Err: err})
		return
	}

	// ------------- Optional query parameter "pageSize" -------------

	err = runtime.BindQueryParameter("form", true, false, "pageSize", r.URL.Query(), &params.PageSize)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "pageSize", Err: err})
		return
	}

	// ------------- Optional query parameter "sortBy" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortBy", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sortBy", Err: err})
		return
	}

	// ------------- Optional query parameter "sortDirection" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortDirection", r.URL.Query(), &params.SortDirection)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sortDirection", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetUserActivity(w, r, shortId, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetInsights operation middleware
func (siw *ServerInterfaceWrapper) GetInsights(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortId" -------------
	var shortId string

	err = runtime.BindStyledParameterWithOptions("simple", "shortId", mux.Vars(r)["shortId"], &shortId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetInsights(w, r, shortId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// ConnectWs operation middleware
func (siw *ServerInterfaceWrapper) ConnectWs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "user" -------------
	var user string

	err = runtime.BindStyledParameterWithOptions("simple", "user", mux.Vars(r)["user"], &user, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "user", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ConnectWs(w, r, user)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{})
}

type GorillaServerOptions struct {
	BaseURL          string
	BaseRouter       *mux.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *mux.Router) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r *mux.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options GorillaServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = mux.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.HandleFunc(options.BaseURL+"/api/v1/bookmarks", wrapper.GetBookmarks).Methods("GET")

	r.HandleFunc(options.BaseURL+"/api/v1/bookmarks", wrapper.CreateBookmark).Methods("POST")

	r.HandleFunc(options.BaseURL+"/api/v1/bookmarks/{shortId}", wrapper.DeleteBookmark).Methods("DELETE")

	r.HandleFunc(options.BaseURL+"/api/v1/bookmarks/{shortId}", wrapper.GetBookmark).Methods("GET")

	r.HandleFunc(options.BaseURL+"/api/v1/bookmarks/{shortId}", wrapper.UpdateBookmark).Methods("PUT")

	r.HandleFunc(options.BaseURL+"/api/v1/notes", wrapper.GetNotes).Methods("GET")

	r.HandleFunc(options.BaseURL+"/api/v1/notes", wrapper.CreateNote).Methods("POST")

	r.HandleFunc(options.BaseURL+"/api/v1/notes/{shortId}", wrapper.DeleteNote).Methods("DELETE")

	r.HandleFunc(options.BaseURL+"/api/v1/notes/{shortId}", wrapper.GetNote).Methods("GET")

	r.HandleFunc(options.BaseURL+"/api/v1/notes/{shortId}", wrapper.UpdateNote).Methods("PUT")

	r.HandleFunc(options.BaseURL+"/api/v1/notes/{shortId}/draft", wrapper.DeleteNoteDraft).Methods("DELETE")

	r.HandleFunc(options.BaseURL+"/api/v1/notes/{shortId}/draft", wrapper.SaveNoteDraft).Methods("PUT")

	r.HandleFunc(options.BaseURL+"/api/v1/tags", wrapper.GetTags).Methods("GET")

	r.HandleFunc(options.BaseURL+"/api/v1/users/{shortId}", wrapper.GetUser).Methods("GET")

	r.HandleFunc(options.BaseURL+"/api/v1/users/{shortId}", wrapper.UpdateUser).Methods("PUT")

	r.HandleFunc(options.BaseURL+"/api/v1/users/{shortId}/activity", wrapper.GetUserActivity).Methods("GET")

	r.HandleFunc(options.BaseURL+"/api/v1/users/{shortId}/insights", wrapper.GetInsights).Methods("GET")

	r.HandleFunc(options.BaseURL+"/api/v1/ws/{user}", wrapper.ConnectWs).Methods("GET")

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"1spA8eAB55JH92gYuPo7AAD//xc9GLp2LQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
