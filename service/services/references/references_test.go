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
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{ShortId: "noteID", CreatedBy: "user"}, nil)
		noteDb.On("GetNotes", []string{"note1", "note2"}, []string{"id", "title"}).Return([]notes.NoteDocument{{ShortId: "note1", Title: "title of note1"}, {ShortId: "note2", Title: "title of note2"}}, nil)
		referencesDb := new(references.MockedReferences)
		referencesDb.On("ListReferencesOfNoteInDepth", "noteID", 1).Return([]references.ReferencesDocument{{From: "note1", To: "note2"}}, nil)
		referencesController := NewService(noteDb, referencesDb)

		// act
		page, error := referencesController.ListReferencesToNote("user", "noteID", 1)

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
		assert.Equal(page, ReferencesResponse{Links: expectedLinks, Nodes: expectedNodes})
	})

	t.Run("it returns error if the getnotes method fails", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{ShortId: "noteID", CreatedBy: "user"}, nil)
		noteDb.On("GetNotes", []string{"note1", "note2"}, []string{"id", "title"}).Return([]notes.NoteDocument{}, fmt.Errorf("sorry"))
		referencesDb := new(references.MockedReferences)
		referencesDb.On("ListReferencesOfNoteInDepth", "noteID", 1).Return([]references.ReferencesDocument{{From: "note1", To: "note2"}}, nil)
		referencesController := NewService(noteDb, referencesDb)

		// act
		page, error := referencesController.ListReferencesToNote("user", "noteID", 1)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("sorry"))
		assert.Equal(page, ReferencesResponse{})
	})

	t.Run("it returns error if the ListReferencesOfNoteInDepth method fails", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{ShortId: "noteID", CreatedBy: "user"}, nil)
		noteDb.On("GetNotes", []string{}, []string{"id", "title"}).Return([]notes.NoteDocument{}, fmt.Errorf("sorry"))
		referencesDb := new(references.MockedReferences)
		referencesDb.On("ListReferencesOfNoteInDepth", "noteID", 1).Return([]references.ReferencesDocument{}, fmt.Errorf("sorry"))
		referencesController := NewService(noteDb, referencesDb)

		// act
		page, error := referencesController.ListReferencesToNote("user", "noteID", 1)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("sorry"))
		assert.Equal(page, ReferencesResponse{})
	})

	t.Run("it returns error if the given note is not found", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{}, fmt.Errorf("error"))
		referencesController := NewService(noteDb, nil)

		// act
		page, error := referencesController.ListReferencesToNote("user", "noteID", 1)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("error"))
		assert.Equal(page, ReferencesResponse{})
	})

	t.Run("it returns error if the given note does not belong to requesting user", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{ShortId: "noteID", CreatedBy: "user1"}, nil)
		referencesController := NewService(noteDb, nil)

		// act
		page, error := referencesController.ListReferencesToNote("user", "noteID", 1)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOT_ALLOWED_TO_GET"))
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
		noteDb.On("GetNotes", []string{"note1", "note2"}, []string{"ID"}).Return([]notes.NoteDocument{}, nil)
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(fmt.Errorf("can not delete"))

		referencesController := NewService(noteDb, referencesDb)

		// act
		error := referencesController.UpsertReferencesOfNote("noteID", noteContent)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("can not delete"))
	})

	t.Run("while upserting references, it returns error if GetNotes fails", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetNotes", []string{"note1", "note2", "note3"}, []string{"id"}).Return([]notes.NoteDocument{}, fmt.Errorf("getNotes failed"))
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(nil)

		referencesController := NewService(noteDb, referencesDb)

		// act
		error := referencesController.UpsertReferencesOfNote("noteID", noteContent)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("getNotes failed"))
	})

	t.Run("while upserting references, it returns error if insertion of references fails", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetNotes", []string{"note1", "note2", "note3"}, []string{"id"}).Return([]notes.NoteDocument{{ShortId: "note1"}, {ShortId: "note3"}}, nil)
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(nil)
		referencesDb.On("InsertMany", "noteID", []string{"note1", "note3"}).Return(fmt.Errorf("insertion failed"))

		referencesController := NewService(noteDb, referencesDb)

		// act
		error := referencesController.UpsertReferencesOfNote("noteID", noteContent)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("insertion failed"))
	})

	t.Run("while upserting references, it only adds references to existing notes", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetNotes", []string{"note1", "note2", "note3"}, []string{"id"}).Return([]notes.NoteDocument{{ShortId: "note1"}, {ShortId: "note3"}}, nil)
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(nil)
		referencesDb.On("InsertMany", "noteID", []string{"note1", "note3"}).Return(nil)

		referencesController := NewService(noteDb, referencesDb)

		// act
		error := referencesController.UpsertReferencesOfNote("noteID", noteContent)

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
	})
}

func TestDeleteReferencesOfNote(t *testing.T) {
	t.Run("it returns error if it can not delete references from the given note", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{}, fmt.Errorf("error"))
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(fmt.Errorf("can not delete"))
		referencesDb.On("DeleteAllReferencesToNote", "noteID").Return(nil)

		referencesController := NewService(noteDb, referencesDb)

		// act
		error := referencesController.DeleteReferencesOfNote("noteID")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("can not delete"))
	})

	t.Run("it returns error if it can not delete references to the given note", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{}, fmt.Errorf("error"))
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(nil)
		referencesDb.On("DeleteAllReferencesToNote", "noteID").Return(fmt.Errorf("can not delete"))

		referencesController := NewService(noteDb, referencesDb)

		// act
		error := referencesController.DeleteReferencesOfNote("noteID")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("can not delete"))
	})

	t.Run("it deletes all references to and from a note, when requested", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{}, fmt.Errorf("error"))
		referencesDb := new(references.MockedReferences)
		referencesDb.On("DeleteAllReferencesFromNote", "noteID").Return(nil)
		referencesDb.On("DeleteAllReferencesToNote", "noteID").Return(nil)

		referencesController := NewService(noteDb, referencesDb)

		// act
		error := referencesController.DeleteReferencesOfNote("noteID")

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
	})
}

/*
	t.Run("it returns an array of history entries and meta information of the page, when list of history entries is requested", func(t *testing.T) {
		// arrange mocks
		paginationMeta := PaginationMeta{
			Count:    0,
			Page:     1,
			PageSize: 10,
			Range:    "",
		}
		currentNote := notes.NoteDocument{ID: "noteID", CreatedBy: "userID", UpdatedBy: "userID", CreatedAt: time.Now(), UpdatedAt: time.Now(), Title: "title", Content: "test", Tags: []string{"test"}}
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(currentNote, nil)
		noteDb.On("List", "userID", "noteID", 1, 10, "title", "asc", []string(nil)).Return([]notes.NoteDocument(nil), 0, nil)

		notesHistoryDb := new(notesHistory.MockedNotesHistory)
		notesHistoryDb.On("ListHistoryOfNote", "userID", "noteID", 1, 10).Return([]notesHistory.NotesHistoryDocument{}, 0, nil)
		notesHistoryController := NewNotesHistoryController(noteDb, notesHistoryDb)

		// act
		listResponse, error := notesHistoryController.ListNotesHistory("userID", "noteID", 1, 10)

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
		assert.Equal(listResponse.Meta, paginationMeta)
	})

	t.Run("it throws an error if parent note of requested history entry does not exist", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteID").Return(notes.NoteDocument{}, fmt.Errorf("error"))

		notesHistoryController := NewNotesHistoryController(noteDb, nil)

		// act
		_, error := notesHistoryController.GetNotesHistory("userID", "noteID", "0001-01-01T22:34:45Z")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOTE_NOT_FOUND"))
	}) */
