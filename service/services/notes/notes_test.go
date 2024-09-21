package notesService

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

func TestNoteService(t *testing.T) {
	t.Run("returns an array of notes and meta information of the page when list of notes is requested", func(t *testing.T) {
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
		page := []notes.NoteDocument{{ShortId: "note1"}, {ShortId: "note2"}}
		hasDraft := map[string]bool{"note1": true}

		stringArrEmpty := []string{}
		noteDb := new(notes.MockedNotes)
		noteDb.On("List", "userId", "test", 1, 10, "title", "asc", stringArrEmpty).Return(page, 0, nil)

		notesDraftsDb := new(notesDrafts.MockedNotesDrafts)
		notesDraftsDb.On("CheckExistence", []string{"note1", "note2"}).Return(hasDraft, nil)

		notesController := NewService(nil, noteDb, nil, notesDraftsDb, nil, nil)

		// act
		listResponse, err := notesController.ListNotes("userId", PointerTo("test"), &paginationMeta.Page, &paginationMeta.PageSize, PointerTo("title"), PointerTo("asc"), &stringArrEmpty)

		// assert
		assert := assert.New(t)
		assert.NoError(err)
		assert.Equal(paginationMeta, listResponse.Meta)
		assert.True(listResponse.Items[0].HasDraft)
		assert.False(listResponse.Items[1].HasDraft)
	})

	t.Run("returns notes and meta information even if CheckExistence fails", func(t *testing.T) {
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
		page := []notes.NoteDocument{{ShortId: "note1"}, {ShortId: "note2"}}

		stringArrEmpty := []string{}
		noteDb := new(notes.MockedNotes)
		noteDb.On("List", "userId", "test", 1, 10, "title", "asc", stringArrEmpty).Return(page, 0, nil)

		draftsResult := make(map[string]bool)
		notesDraftsDb := new(notesDrafts.MockedNotesDrafts)
		notesDraftsDb.On("CheckExistence", []string{"note1", "note2"}).Return(draftsResult, fmt.Errorf("failed"))

		notesController := NewService(nil, noteDb, nil, notesDraftsDb, nil, nil)

		// act
		listResponse, err := notesController.ListNotes("userId", PointerTo("test"), &paginationMeta.Page, &paginationMeta.PageSize, PointerTo("title"), PointerTo("asc"), &stringArrEmpty)

		// assert
		assert := assert.New(t)
		assert.NoError(err)
		assert.Equal(paginationMeta, listResponse.Meta)
		assert.False(listResponse.Items[0].HasDraft)
		assert.False(listResponse.Items[1].HasDraft)
	})

	t.Run("throws an error while updating a non-existing note", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{}, nil)

		notesController := NewService(nil, noteDb, nil, nil, nil, nil)

		// act
		err := notesController.UpdateNote("noteId", "user2", "clientId", nil)

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "NOT_ALLOWED_TO_UPDATE")
	})

	t.Run("throws an error while deleting a non-existing note", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{}, nil)

		notesController := NewService(nil, noteDb, nil, nil, nil, nil)

		// act
		err := notesController.DeleteNote("noteId", "user1", "clientId")

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "NOT_ALLOWED_TO_DELETE")
	})

	t.Run("does not allow getting a note if not created by the requester", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{CreatedBy: "user1", Id: "Id"}, nil)

		notesController := NewService(nil, noteDb, nil, nil, nil, nil)

		// act
		_, err := notesController.GetNote("noteId", "user2", "clientId", GetNoteParams{})

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "NOT_ALLOWED_TO_GET")
	})

	t.Run("does not allow updating a note not created by the requester", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{CreatedBy: "user1", Id: "Id"}, nil)

		notesController := NewService(nil, noteDb, nil, nil, nil, nil)

		// act
		err := notesController.UpdateNote("noteId", "user2", "", nil)

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "NOT_ALLOWED_TO_UPDATE")
	})

	t.Run("does not allow deleting a note not created by the requester", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{CreatedBy: "user1", Id: "Id"}, nil)

		notesController := NewService(nil, noteDb, nil, nil, nil, nil)

		// act
		err := notesController.DeleteNote("noteId", "user2", "")

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "NOT_ALLOWED_TO_DELETE")
	})

	t.Run("throws an error while querying the draft of a non-existing note", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{}, fmt.Errorf("err"))

		notesController := NewService(nil, noteDb, nil, nil, nil, nil)

		// act
		_, err := notesController.GetDraftOfNote("userId", "noteId")

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "err")
	})

	t.Run("throws an error while updating the draft of a note belonging to another user", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{Id: "noteId", CreatedBy: "user0"}, nil)

		notesController := NewService(nil, noteDb, nil, nil, nil, nil)

		// act
		err := notesController.UpdateDraft("user1", "noteId", "title", "", nil)

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "NOT_ALLOWED_TO_UPDATE")
	})

	t.Run("throws an error while deleting the draft of a note belonging to another user", func(t *testing.T) {
		// arrange mocks
		noteDb := new(notes.MockedNotes)
		noteDb.On("GetOne", "noteId").Return(notes.NoteDocument{Id: "noteId", CreatedBy: "user0"}, nil)

		notesController := NewService(nil, noteDb, nil, nil, nil, nil)

		// act
		err := notesController.DeleteDraft("user1", "noteId")

		// assert
		assert := assert.New(t)
		assert.EqualError(err, "NOT_ALLOWED_TO_DELETE")
	})
}
