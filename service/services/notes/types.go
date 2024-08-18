package notesService

import (
	"time"

	referencesService "github.com/alperklc/the-zula/service/services/references"
)

type GetNoteParams struct {
	LoadDraft     bool
	GetHistory    bool
	GetReferences bool
}

type Note struct {
	ShortId    string                                `json:"id"`
	UpdatedAt  time.Time                             `json:"updatedAt"`
	UpdatedBy  string                                `json:"updatedBy"`
	CreatedBy  string                                `json:"createdBy"`
	CreatedAt  time.Time                             `json:"createdAt"`
	Title      string                                `json:"title"`
	Content    string                                `json:"content"`
	Tags       []string                              `json:"tags"`
	HasDraft   bool                                  `json:"hasDraft"`
	Versions   *int32                                `json:"versions,omitempty"`
	References *referencesService.ReferencesResponse `json:"references,omitempty"`
}

type PaginationMeta struct {
	Count         int    `json:"count"`
	Query         string `json:"query"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
	SortBy        string `json:"sortBy"`
	SortDirection string `json:"sortDirection"`
	Range         string `json:"range"`
}

type NotesPage struct {
	Meta  PaginationMeta `json:"meta"`
	Items []Note         `json:"items"`
}

type Statistics struct {
	Count int64 `json:"count"`
}
