package notesReferences

import (
	"github.com/stretchr/testify/mock"
)

type MockedNotesReferences struct {
	mock.Mock
}

func (m *MockedNotesReferences) ListReferencesOfNoteInDepth(noteId string, depth int) ([]ReferencesDocument, error) {
	args := m.Called(noteId, depth)

	return args.Get(0).([]ReferencesDocument), args.Error(1)
}

func (m *MockedNotesReferences) InsertMany(from string, to []string) error {
	args := m.Called(from, to)

	return args.Error(0)
}

func (m *MockedNotesReferences) DeleteAllReferencesFromNote(noteId string) error {
	args := m.Called(noteId)

	return args.Error(0)
}

func (m *MockedNotesReferences) DeleteAllReferencesToNote(noteId string) error {
	args := m.Called(noteId)

	return args.Error(0)
}
