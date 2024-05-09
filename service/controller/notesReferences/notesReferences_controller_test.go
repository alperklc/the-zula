package notesReferencesCtrl

import (
	"fmt"
	"testing"

	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesReferences"

	"github.com/stretchr/testify/assert"
)

func TestListReferencesToNote(t *testing.T) {
	t.Run("it returns the references of node", func(t *testing.T) {
		// arrange mocks
		expectedLinks := []NoteReferenceLink{{Source: "note1", Target: "note2"}}
		expectedNodes := []NoteReferenceNode{{ID: "note1", Title: "title of note1"}, {ID: "note2", Title: "title of note2"}}
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetOne", "noteUID").Return(notes.NoteDocument{UID: "noteUID", CreatedBy: "user"}, nil)
		noteRepository.On("GetNotes", []string{"note1", "note2"}, []string{"uid", "title"}).Return([]notes.NoteDocument{{UID: "note1", Title: "title of note1"}, {UID: "note2", Title: "title of note2"}}, nil)
		notesReferencesRepository := new(notesReferences.MockedNotesReferences)
		notesReferencesRepository.On("ListReferencesOfNoteInDepth", "noteUID", 1).Return([]notesReferences.ReferencesDocument{{From: "note1", To: "note2"}}, nil)
		notesReferencesController := NewNotesReferencesController(noteRepository, notesReferencesRepository)

		// act
		page, error := notesReferencesController.ListReferencesToNote("user", "noteUID", 1)

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
		assert.Equal(page, NoteReferencesResponse{Links: expectedLinks, Nodes: expectedNodes})
	})

	t.Run("it returns error if the getnotes method fails", func(t *testing.T) {
		// arrange mocks
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetOne", "noteUID").Return(notes.NoteDocument{UID: "noteUID", CreatedBy: "user"}, nil)
		noteRepository.On("GetNotes", []string{"note1", "note2"}, []string{"uid", "title"}).Return([]notes.NoteDocument{}, fmt.Errorf("sorry"))
		notesReferencesRepository := new(notesReferences.MockedNotesReferences)
		notesReferencesRepository.On("ListReferencesOfNoteInDepth", "noteUID", 1).Return([]notesReferences.ReferencesDocument{{From: "note1", To: "note2"}}, nil)
		notesReferencesController := NewNotesReferencesController(noteRepository, notesReferencesRepository)

		// act
		page, error := notesReferencesController.ListReferencesToNote("user", "noteUID", 1)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("sorry"))
		assert.Equal(page, NoteReferencesResponse{})
	})

	t.Run("it returns error if the ListReferencesOfNoteInDepth method fails", func(t *testing.T) {
		// arrange mocks
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetOne", "noteUID").Return(notes.NoteDocument{UID: "noteUID", CreatedBy: "user"}, nil)
		noteRepository.On("GetNotes", []string{}, []string{"uid", "title"}).Return([]notes.NoteDocument{}, fmt.Errorf("sorry"))
		notesReferencesRepository := new(notesReferences.MockedNotesReferences)
		notesReferencesRepository.On("ListReferencesOfNoteInDepth", "noteUID", 1).Return([]notesReferences.ReferencesDocument{}, fmt.Errorf("sorry"))
		notesReferencesController := NewNotesReferencesController(noteRepository, notesReferencesRepository)

		// act
		page, error := notesReferencesController.ListReferencesToNote("user", "noteUID", 1)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("sorry"))
		assert.Equal(page, NoteReferencesResponse{})
	})

	t.Run("it returns error if the given note is not found", func(t *testing.T) {
		// arrange mocks
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetOne", "noteUID").Return(notes.NoteDocument{}, fmt.Errorf("error"))
		notesReferencesController := NewNotesReferencesController(noteRepository, nil)

		// act
		page, error := notesReferencesController.ListReferencesToNote("user", "noteUID", 1)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOTE_NOT_FOUND"))
		assert.Equal(page, NoteReferencesResponse{})
	})

	t.Run("it returns error if the given note does not belong to requesting user", func(t *testing.T) {
		// arrange mocks
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetOne", "noteUID").Return(notes.NoteDocument{UID: "noteUID", CreatedBy: "user1"}, nil)
		notesReferencesController := NewNotesReferencesController(noteRepository, nil)

		// act
		page, error := notesReferencesController.ListReferencesToNote("user", "noteUID", 1)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOT_ALLOWED_TO_GET"))
		assert.Equal(page, NoteReferencesResponse{})
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
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetNotes", []string{"note1", "note2"}, []string{"uid"}).Return([]notes.NoteDocument{}, nil)
		notesReferencesRepository := new(notesReferences.MockedNotesReferences)
		notesReferencesRepository.On("DeleteAllReferencesFromNote", "noteUID").Return(fmt.Errorf("can not delete"))

		notesReferencesController := NewNotesReferencesController(noteRepository, notesReferencesRepository)

		// act
		error := notesReferencesController.UpsertReferencesOfNote("noteUID", noteContent)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("can not delete"))
	})

	t.Run("while upserting references, it returns error if GetNotes fails", func(t *testing.T) {
		// arrange mocks
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetNotes", []string{"note1", "note2", "note3"}, []string{"uid"}).Return([]notes.NoteDocument{}, fmt.Errorf("getNotes failed"))
		notesReferencesRepository := new(notesReferences.MockedNotesReferences)
		notesReferencesRepository.On("DeleteAllReferencesFromNote", "noteUID").Return(nil)

		notesReferencesController := NewNotesReferencesController(noteRepository, notesReferencesRepository)

		// act
		error := notesReferencesController.UpsertReferencesOfNote("noteUID", noteContent)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("getNotes failed"))
	})

	t.Run("while upserting references, it returns error if insertion of references fails", func(t *testing.T) {
		// arrange mocks
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetNotes", []string{"note1", "note2", "note3"}, []string{"uid"}).Return([]notes.NoteDocument{{UID: "note1"}, {UID: "note3"}}, nil)
		notesReferencesRepository := new(notesReferences.MockedNotesReferences)
		notesReferencesRepository.On("DeleteAllReferencesFromNote", "noteUID").Return(nil)
		notesReferencesRepository.On("InsertMany", "noteUID", []string{"note1", "note3"}).Return(fmt.Errorf("insertion failed"))

		notesReferencesController := NewNotesReferencesController(noteRepository, notesReferencesRepository)

		// act
		error := notesReferencesController.UpsertReferencesOfNote("noteUID", noteContent)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("insertion failed"))
	})

	t.Run("while upserting references, it only adds references to existing notes", func(t *testing.T) {
		// arrange mocks
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetNotes", []string{"note1", "note2", "note3"}, []string{"uid"}).Return([]notes.NoteDocument{{UID: "note1"}, {UID: "note3"}}, nil)
		notesReferencesRepository := new(notesReferences.MockedNotesReferences)
		notesReferencesRepository.On("DeleteAllReferencesFromNote", "noteUID").Return(nil)
		notesReferencesRepository.On("InsertMany", "noteUID", []string{"note1", "note3"}).Return(nil)

		notesReferencesController := NewNotesReferencesController(noteRepository, notesReferencesRepository)

		// act
		error := notesReferencesController.UpsertReferencesOfNote("noteUID", noteContent)

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
	})
}

func TestDeleteReferencesOfNote(t *testing.T) {
	t.Run("it returns error if it can not delete references from the given note", func(t *testing.T) {
		// arrange mocks
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetOne", "noteUID").Return(notes.NoteDocument{}, fmt.Errorf("error"))
		notesReferencesRepository := new(notesReferences.MockedNotesReferences)
		notesReferencesRepository.On("DeleteAllReferencesFromNote", "noteUID").Return(fmt.Errorf("can not delete"))
		notesReferencesRepository.On("DeleteAllReferencesToNote", "noteUID").Return(nil)

		notesReferencesController := NewNotesReferencesController(noteRepository, notesReferencesRepository)

		// act
		error := notesReferencesController.DeleteReferencesOfNote("noteUID")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("can not delete"))
	})

	t.Run("it returns error if it can not delete references to the given note", func(t *testing.T) {
		// arrange mocks
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetOne", "noteUID").Return(notes.NoteDocument{}, fmt.Errorf("error"))
		notesReferencesRepository := new(notesReferences.MockedNotesReferences)
		notesReferencesRepository.On("DeleteAllReferencesFromNote", "noteUID").Return(nil)
		notesReferencesRepository.On("DeleteAllReferencesToNote", "noteUID").Return(fmt.Errorf("can not delete"))

		notesReferencesController := NewNotesReferencesController(noteRepository, notesReferencesRepository)

		// act
		error := notesReferencesController.DeleteReferencesOfNote("noteUID")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("can not delete"))
	})

	t.Run("it deletes all references to and from a note, when requested", func(t *testing.T) {
		// arrange mocks
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetOne", "noteUID").Return(notes.NoteDocument{}, fmt.Errorf("error"))
		notesReferencesRepository := new(notesReferences.MockedNotesReferences)
		notesReferencesRepository.On("DeleteAllReferencesFromNote", "noteUID").Return(nil)
		notesReferencesRepository.On("DeleteAllReferencesToNote", "noteUID").Return(nil)

		notesReferencesController := NewNotesReferencesController(noteRepository, notesReferencesRepository)

		// act
		error := notesReferencesController.DeleteReferencesOfNote("noteUID")

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
		currentNote := notes.NoteDocument{UID: "noteUID", CreatedBy: "userUID", UpdatedBy: "userUID", CreatedAt: time.Now(), UpdatedAt: time.Now(), Title: "title", Content: "test", Tags: []string{"test"}}
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetOne", "noteUID").Return(currentNote, nil)
		noteRepository.On("List", "userUID", "noteUID", 1, 10, "title", "asc", []string(nil)).Return([]notes.NoteDocument(nil), 0, nil)

		notesHistoryRepository := new(notesHistory.MockedNotesHistory)
		notesHistoryRepository.On("ListHistoryOfNote", "userUID", "noteUID", 1, 10).Return([]notesHistory.NotesHistoryDocument{}, 0, nil)
		notesHistoryController := NewNotesHistoryController(noteRepository, notesHistoryRepository)

		// act
		listResponse, error := notesHistoryController.ListNotesHistory("userUID", "noteUID", 1, 10)

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
		assert.Equal(listResponse.Meta, paginationMeta)
	})

	t.Run("it throws an error if parent note of requested history entry does not exist", func(t *testing.T) {
		// arrange mocks
		noteRepository := new(notes.MockedNotes)
		noteRepository.On("GetOne", "noteUID").Return(notes.NoteDocument{}, fmt.Errorf("error"))

		notesHistoryController := NewNotesHistoryController(noteRepository, nil)

		// act
		_, error := notesHistoryController.GetNotesHistory("userUID", "noteUID", "0001-01-01T22:34:45Z")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOTE_NOT_FOUND"))
	}) */
