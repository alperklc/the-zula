package notes

import (
	"time"

	"github.com/alperklc/the-zula/service/utils"
)

type NoteDocument struct {
	Id        string    `bson:"_id,omitempty"`
	ShortId   string    `bson:"shortId"`
	UpdatedAt time.Time `bson:"updatedAt"`
	UpdatedBy string    `bson:"updatedBy"`
	CreatedAt time.Time `bson:"createdAt"`
	CreatedBy string    `bson:"createdBy"`
	Title     string    `bson:"title"`
	Content   string    `bson:"content"`
	Tags      []string  `bson:"tags"`
}

func (n NoteDocument) GetId() string {
	return n.Id
}

func (nd *NoteDocument) IsDifferent(otherNote NoteDocument) bool {
	titleChanged := otherNote.Title != nd.Title
	contentChanged := otherNote.Content != nd.Content
	tagsChanged := !utils.AreArraysEqual(otherNote.Tags, nd.Tags)

	return titleChanged || contentChanged || tagsChanged
}

type NoteFacets struct {
	Count int `json:"count"`
}

type NotesAggregation struct {
	Facets NoteFacets     `json:"facets"`
	Items  []NoteDocument `json:"items"`
}

type TagsResult struct {
	Value string `bson:"_id" json:"value,omitempty"`
	Count int    `bson:"count" json:"count"`
}
