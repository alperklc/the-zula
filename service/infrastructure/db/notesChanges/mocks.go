package notesChanges

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockedNotesChanges struct {
	mock.Mock
}

func (m *MockedNotesChanges) ListHistoryOfNote(userId, noteId string, page, pageSize int) ([]NotesChangesDocument, int, error) {
	args := m.Called(userId, noteId, page, pageSize)

	return args.Get(0).([]NotesChangesDocument), args.Int(1), args.Error(2)
}

func (m *MockedNotesChanges) GetCountOfChanges(noteId string) (int64, error) {
	args := m.Called(noteId)

	return args.Get(0).(int64), args.Error(1)
}

func (m *MockedNotesChanges) GetOne(shortId string) (NotesChangesDocument, error) {
	args := m.Called(shortId)

	return args.Get(0).(NotesChangesDocument), args.Error(1)
}

func (m *MockedNotesChanges) InsertOne(noteId string, updatedAt time.Time, updatedBy, change string) error {
	args := m.Called(noteId, updatedAt, updatedBy, change)

	return args.Error(0)
}
