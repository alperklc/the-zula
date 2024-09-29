package notesChanges

import "time"

type NotesChangesDocument struct {
	Id        string    `bson:"_id,omitempty"`
	ShortId   string    `bson:"shortId"`
	NoteId    string    `bson:"noteId"`
	UpdatedAt time.Time `bson:"updatedAt"`
	UpdatedBy string    `bson:"updatedBy"`
	Change    string    `bson:"change"`
}

type NotesChangesFacets struct {
	Count int `json:"count"`
}

type NotesChangesAggregation struct {
	Facets NotesChangesFacets     `json:"facets"`
	Items  []NotesChangesDocument `json:"items"`
}

func (n NotesChangesDocument) GetId() string {
	return n.Id
}
