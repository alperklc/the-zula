package notectrl

import (
	"fmt"
	"testing"

	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesDrafts"

	"github.com/stretchr/testify/assert"
)

func PointerTo[T ~string](s T) *T {
	return &s
}

func TestNote(t *testing.T) {
	t.Run("it returns an array of notes and meta information of the page, when list of notes is requested", func(t *testing.T) {
		// arrange mocks
		paginationMeta := PaginationMeta{
			Count:         0,
			Query:         "test",
			Page:          1,
			PageSize:      10,
			SortBy:        "title",
			SortDirection: "asc",
			Range:         "",
		}
		page := []notes.NoteDocument{{Id: "note1"}, {Id: "note2"}}
		hasDraft := make(map[string]bool)
		hasDraft["note1"] = true

		noteDb := new(notes.MockedNotes)
		noteDb.On("List", "userId", "test", 1, 10, "title", "asc", []string(nil)).Return(page, 0, nil)
		notesDraftsDb := new(notesDrafts.MockedNotesDrafts)
		notesDraftsDb.On("CheckExistence", "userId", []string{"note1", "note2"}).Return(hasDraft, nil)

		notesController := NewNotesController(noteDb, notesDraftsDb)

		// act
		listResponse, error := notesController.ListNotes("userId", PointerTo("test"), 1, 10, PointerTo("title"), PointerTo("asc"), []string(nil))

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
		assert.Equal(listResponse.Meta, paginationMeta)
		assert.True(listResponse.Items[0].HasDraft)
		assert.False(listResponse.Items[1].HasDraft)
	})

	t.Run("it returns an array of notes and meta information of the page, even when checkexistence fails", func(t *testing.T) {
		// arrange mocks
		paginationMeta := PaginationMeta{
			Count:         0,
			Query:         "test",
			Page:          1,
			PageSize:      10,
			SortBy:        "title",
			SortDirection: "asc",
			Range:         "",
		}
		page := []notes.NoteDocument{{Id: "note1"}, {Id: "note2"}}
		hasDraft := make(map[string]bool)

		noteDb := new(notes.MockedNotes)
		noteDb.On("List", "userId", "test", 1, 10, "title", "asc", []string(nil)).Return(page, 0, nil)
		notesDraftsDb := new(notesDrafts.MockedNotesDrafts)
		notesDraftsDb.On("CheckExistence", "userId", []string{"note1", "note2"}).Return(hasDraft, fmt.Errorf("failed"))

		notesController := NewNotesController(noteDb, notesDraftsDb)

		// act
		listResponse, error := notesController.ListNotes("userId", "test", 1, 10, "title", "asc", []string(nil))

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
		assert.Equal(listResponse.Meta, paginationMeta)
		assert.False(listResponse.Items[0].HasDraft)
		assert.False(listResponse.Items[1].HasDraft)
	})
	/*
		t.Run("it allows only certain fields while updating a note and deletes draft entry after updating", func(t *testing.T) {
			// arrange mocks
			initialUpdates := make(map[string]interface{})
			initialUpdates["foo"] = "bar"
			initialUpdates["content"] = "test"

			expectedUpdates := make(map[string]interface{})
			expectedUpdates["content"] = "test"

			currentNote := notes.NoteDocument{Id: "noteId", CreatedBy: "userId", UpdatedBy: "userId", CreatedAt: time.Now(), UpdatedAt: time.Now(), Title: "title", Content: "old", Tags: []string{"test"}}
			noteMessage := messageQueue.NoteMessage{UserId: "userId", Action: "UPDATE", ObjectId: "noteId", ClientID: "clientId"}
			noteDb := new(notes.MockedNotes)
			noteDb.On("GetOne", "noteId").Return(currentNote, nil)
			noteDb.On("UpdateOne", "userId", "noteId", expectedUpdates).Return(nil)
			notesDraftsDb := new(notesDrafts.MockedNotesDrafts)
			notesDraftsDb.On("DeleteOne", "userId", "noteId").Return(nil)

			eventPublisher := new(noteMessagePublisher.MockedEventPublisher)
			eventPublisher.On("Publish", noteMessage).Return(nil)
			notesController := NewNotesController(noteDb, notesDraftsDb)

			// act
			error := notesController.UpdateNote("noteId", "userId", "clientId", initialUpdates)

			// assert
			assert := assert.New(t)
			assert.Equal(error, nil)
			noteDb.MethodCalled("UpdateOne", "userId", "noteId", expectedUpdates)
			notesDraftsDb.MethodCalled("DeleteOne", "userId", "noteId")
		})
	*/
	t.Run("it throws an error while updating a non existing note", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{}, nil)
		notesController := NewNotesController(noteDb, nil)

		// act
		error := notesController.UpdateNote("noteId", "user2", "clientId", nil)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOT_FOUND"))
	})

	t.Run("it throws an error while deleting a non existing note", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{}, nil)
		notesController := NewNotesController(noteDb, nil)

		// act
		error := notesController.DeleteNote("noteId", "user1", "clientId")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOT_FOUND"))
	})
	/*
		t.Run("it publishes a message to the queue if a note is read", func(t *testing.T) {
			// arrange mocks
			noteMessage := messageQueue.NoteMessage{UserId: "userId", Action: "READ", ObjectId: "noteId", ClientID: "clientId"}

			noteDb := new(notes.MockedNotes)
			noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{CreatedBy: "userId", Id: "noteId"}, nil)
			eventPublisher := new(noteMessagePublisher.MockedEventPublisher)
			eventPublisher.On("Publish", noteMessage).Return(nil)
			notesDraftsDb := new(notesDrafts.MockedNotesDrafts)
			notesDraftsDb.On("CheckExistence", "userId", []string{"noteId"}).Return(make(map[string]bool), nil)

			notesController := NewNotesController(noteDb, nil, notesDraftsDb, eventPublisher)

			// act
			_, error := notesController.GetNote("noteId", "userId", "clientId", true)

			// assert
			assert := assert.New(t)
			assert.Equal(error, nil)
			eventPublisher.AssertCalled(t, "Publish", noteMessage)
		})

		t.Run("it won't publish a message to the queue, if a note is read and tracking is opted out", func(t *testing.T) {
			// arrange mocks
			noteMessage := messageQueue.NoteMessage{UserId: "userId", Action: "READ", ObjectId: "noteId", ClientID: "clientId"}

			noteDb := new(notes.MockedNotes)
			noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{CreatedBy: "userId", Id: "noteId"}, nil)
			eventPublisher := new(noteMessagePublisher.MockedEventPublisher)
			eventPublisher.On("Publish", noteMessage).Return(nil)
			notesDraftsDb := new(notesDrafts.MockedNotesDrafts)
			notesDraftsDb.On("CheckExistence", "userId", []string{"noteId"}).Return(make(map[string]bool), nil)

			notesController := NewNotesController(noteDb, nil, notesDraftsDb, eventPublisher)

			// act
			_, error := notesController.GetNote("noteId", "userId", "clientId", false)

			// assert
			assert := assert.New(t)
			assert.Equal(error, nil)
			eventPublisher.AssertNotCalled(t, "Publish", noteMessage)
		})
	*/
	t.Run("it wont allow getting a note if the note is not created by the requester", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{CreatedBy: "user1", Id: "Id"}, nil)
		notesController := NewNotesController(noteDb, nil)

		// act
		_, error := notesController.GetNote("noteId", "user2", "clientId")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOT_ALLOWED_TO_GET"))
	})

	t.Run("it returns the notes only if they belong to requester, if notes requested within a bulk request", func(t *testing.T) {
		// arrange mocks
		noteIds := []string{"note1", "note2", "note3"}
		fields := []string{"title", "createdBy"}
		note1 := notes.NoteDocument{CreatedBy: "user", Id: "note1"}
		note2 := notes.NoteDocument{CreatedBy: "user2", Id: "note2"}
		note3 := notes.NoteDocument{CreatedBy: "user", Id: "note3"}
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetNotes", noteIds, fields).Return([]notes.NoteDocument{note1, note2, note3}, nil)
		notesController := NewNotesController(noteDb, nil)

		// act
		notesReturned, error := notesController.GetNotes("user", noteIds, fields)

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
		assert.Contains(notesReturned, Note{Id: note1.Id, CreatedBy: note1.CreatedBy})
		assert.Contains(notesReturned, Note{Id: note3.Id, CreatedBy: note3.CreatedBy})
		assert.NotContains(notesReturned, Note{Id: note2.Id, CreatedBy: note2.CreatedBy})
	})

	t.Run("it wont allow updating a note if the note is not created by the requester", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{CreatedBy: "user1", Id: "Id"}, nil)
		notesController := NewNotesController(noteDb, nil)

		// act
		error := notesController.UpdateNote("noteId", "user2", "", nil)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOT_ALLOWED_TO_UPDATE"))
	})

	t.Run("it wont allow deleting a note if the note is not created by requester", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{CreatedBy: "user1", Id: "Id"}, nil)
		notesController := NewNotesController(noteDb, nil)

		// act
		error := notesController.DeleteNote("noteId", "user2", "")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOT_ALLOWED_TO_DELETE"))
	})

	/*	t.Run("it deletes draft after a note gets deleted", func(t *testing.T) {
			// arrange mocks
			noteDb := new(notes.MockedNotes)
			noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{CreatedBy: "user1", Id: "noteId"}, nil)
			noteDb.On("DeleteOne", "noteId").Return(nil)
			notesDraftsDb := new(notesDrafts.MockedNotesDrafts)
			notesDraftsDb.On("DeleteOne", "user1", "noteId").Return(nil)
			eventPublisher := new(noteMessagePublisher.MockedEventPublisher)
			noteMessage := messageQueue.NoteMessage{UserId: "user1", Action: "DELETE", ObjectId: "noteId", ClientID: "clientId"}
			eventPublisher.On("Publish", noteMessage).Return(nil)

			notesController := NewNotesController(noteDb, nil, notesDraftsDb, eventPublisher)

			// act
			error := notesController.DeleteNote("noteId", "user1", "clientId")

			// assert
			assert := assert.New(t)
			assert.Equal(error, nil)
			notesDraftsDb.AssertCalled(t, "DeleteOne", "user1", "noteId")
			eventPublisher.AssertCalled(t, "Publish", noteMessage)
		})
	*/
	t.Run("it throws an error while querying draft of a non existing note", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{}, fmt.Errorf("err"))
		notesController := NewNotesController(noteDb, nil)

		// act
		_, error := notesController.GetDraftOfNote("userId", "noteId")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOT_FOUND"))
	})

	t.Run("it throws an error while updating draft of a note belongs to another user", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{Id: "noteId", CreatedBy: "user0"}, nil)
		notesController := NewNotesController(noteDb, nil)

		// act
		error := notesController.UpdateDraft("user1", "noteId", "title", "", nil)

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOT_ALLOWED_TO_UPDATE"))
	})

	t.Run("it throws an error while updating draft of a note belongs to another user", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{Id: "noteId", CreatedBy: "user0"}, nil)
		notesController := NewNotesController(noteDb, nil)

		// act
		error := notesController.DeleteDraft("user1", "noteId")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("NOT_ALLOWED_TO_DELETE"))
	})
}
