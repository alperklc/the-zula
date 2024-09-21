package referencesService

import (
	"fmt"
	"testing"

	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/references"

	"github.com/stretchr/testify/assert"
)

func TestListReferencesToNote(t *testing.T) {
	t.Run("it returns the references of node", func(t *testing.T) {
		// arrange mocks
		expectedLinks := []ReferenceLink{{Source: "note1", Target: "note2"}}
		expectedNodes := []ReferenceNode{{ID: "note1", Title: "title of note1"}, {ID: "note2", Title: "title of note2"}}
		noteDb := new(notes.MockedNotes) // Updated for interface
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{ShortId: "noteID", CreatedBy: "user"}, nil)
		noteDb.On("GetNotes", []string{"note1", "note2"}, []string{"id", "shortId", "title"}).Return([]notes.NoteDocument{{ShortId: "note1", Title: "title of note1"}, {ShortId: "note2", Title: "title of note2"}}, nil)
		referencesDb := new(references.MockedReferences) // Updated for interface
		referencesDb.On("ListReferencesOfNoteInDepth", "noteID", 1).Return([]references.ReferencesDocument{{From: "note1", To: "note2"}}, nil)
		referencesController := NewService(noteDb, referencesDb)

		// act
		page, err := referencesController.ListReferencesToNote("user", "noteID", 1)

		// assert
		assert := assert.New(t)
		assert.Nil(err)
		assert.Equal(page, ReferencesResponse{Links: expectedLinks, Nodes: expectedNodes})
	})

	t.Run("it returns error if the getnotes method fails", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{ShortId: "noteID", CreatedBy: "user"}, nil)
		noteDb.On("GetNotes", []string{"note1", "note2"}, []string{"id", "shortId", "title"}).Return([]notes.NoteDocument{}, fmt.Errorf("sorry"))
		referencesDb := new(references.MockedReferences)
		referencesDb.On("ListReferencesOfNoteInDepth", "noteID", 1).Return([]references.ReferencesDocument{{From: "note1", To: "note2"}}, nil)
		referencesController := NewService(noteDb, referencesDb)

		// act
		page, err := referencesController.ListReferencesToNote("user", "noteID", 1)

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "sorry")
		assert.Equal(page, ReferencesResponse{})
	})

	t.Run("it returns error if the ListReferencesOfNoteInDepth method fails", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{ShortId: "noteID", CreatedBy: "user"}, nil)
		noteDb.On("GetNotes", []string{}, []string{"id", "shortId", "title"}).Return([]notes.NoteDocument{}, nil)
		referencesDb := new(references.MockedReferences)
		referencesDb.On("ListReferencesOfNoteInDepth", "noteID", 1).Return([]references.ReferencesDocument{}, fmt.Errorf("sorry"))
		referencesController := NewService(noteDb, referencesDb)

		// act
		page, err := referencesController.ListReferencesToNote("user", "noteID", 1)

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "sorry")
		assert.Equal(page, ReferencesResponse{})
	})

	t.Run("it returns error if the given note is not found", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{}, fmt.Errorf("error"))
		referencesController := NewService(noteDb, nil)

		// act
		page, err := referencesController.ListReferencesToNote("user", "noteID", 1)

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "error")
		assert.Equal(page, ReferencesResponse{})
	})

	t.Run("it returns error if the given note does not belong to requesting user", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{ShortId: "noteID", CreatedBy: "user1"}, nil)
		referencesController := NewService(noteDb, nil)

		// act
		page, err := referencesController.ListReferencesToNote("user", "noteID", 1)

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "NOT_ALLOWED_TO_GET")
		assert.Equal(page, ReferencesResponse{})
	})
}

func TestUpsertReferencesOfNote(t *testing.T) {
	noteContent := `this is an example markdown text which contains links to other notes
				  - [first link](/notes/note1)  - this one exist
				  - [second link](/notes/note2) - this one doesnt exist
				  - [third link](/notes/note3)  - this one exist
				`

	t.Run("while upserting references, it returns error if it can not delete current references from the given note", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetNotes", []string{"note1", "note2"}, []string{"id"}).Return([]notes.NoteDocument{}, nil)
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(fmt.Errorf("can not delete"))

		referencesController := NewService(noteDb, referencesDb)

		// act
		err := referencesController.UpsertReferencesOfNote("noteID", noteContent)

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "can not delete")
	})

	t.Run("while upserting references, it returns error if GetNotes fails", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetNotes", []string{"note1", "note2", "note3"}, []string{"id", "shortId"}).Return([]notes.NoteDocument{}, fmt.Errorf("getNotes failed"))
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(nil)

		referencesController := NewService(noteDb, referencesDb)

		// act
		err := referencesController.UpsertReferencesOfNote("noteID", noteContent)

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "getNotes failed")
	})

	t.Run("while upserting references, it returns error if insertion of references fails", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetNotes", []string{"note1", "note2", "note3"}, []string{"id", "shortId"}).Return([]notes.NoteDocument{{ShortId: "note1"}, {ShortId: "note3"}}, nil)
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(nil)
		referencesDb.On("InsertMany", "noteID", []string{"note1", "note3"}).Return(fmt.Errorf("insertion failed"))

		referencesController := NewService(noteDb, referencesDb)

		// act
		err := referencesController.UpsertReferencesOfNote("noteID", noteContent)

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "insertion failed")
	})

	t.Run("while upserting references, it only adds references to existing notes", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetNotes", []string{"note1", "note2", "note3"}, []string{"id", "shortId"}).Return([]notes.NoteDocument{{ShortId: "note1"}, {ShortId: "note3"}}, nil)
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(nil)
		referencesDb.On("InsertMany", "noteID", []string{"note1", "note3"}).Return(nil)

		referencesController := NewService(noteDb, referencesDb)

		// act
		err := referencesController.UpsertReferencesOfNote("noteID", noteContent)

		// assert
		assert := assert.New(t)
		assert.Nil(err)
	})
}

func TestDeleteReferencesOfNote(t *testing.T) {
	t.Run("it returns error if it can not delete references of given note", func(t *testing.T) {
		// arrange mocks
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(fmt.Errorf("deletion failed"))

		referencesController := NewService(nil, referencesDb)

		// act
		err := referencesController.DeleteReferencesOfNote("noteID")

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "deletion failed")
	})

	t.Run("it deletes all the references of given note", func(t *testing.T) {
		// arrange mocks
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(nil)
		referencesDb.On("DeleteAllReferencesToNote", "noteID").Return(nil)

		referencesController := NewService(nil, referencesDb)

		// act
		err := referencesController.DeleteReferencesOfNote("noteID")

		// assert
		assert := assert.New(t)
		assert.Nil(err)
	})
}
