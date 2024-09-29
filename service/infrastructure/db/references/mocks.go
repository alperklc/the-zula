package references

import (
	"github.com/stretchr/testify/mock"
)

type MockedReferences struct {
	mock.Mock
}

func (m *MockedReferences) ListReferencesOfNoteInDepth(noteId string, depth int) ([]ReferencesDocument, error) {
	args := m.Called(noteId, depth)

	return args.Get(0).([]ReferencesDocument), args.Error(1)
}

func (m *MockedReferences) InsertMany(from string, to []string) error {
	args := m.Called(from, to)

	return args.Error(0)
}

func (m *MockedReferences) DeleteAllReferencesFromNote(noteId string) error {
	args := m.Called(noteId)

	return args.Error(0)
}

func (m *MockedReferences) DeleteAllReferencesToNote(noteId string) error {
	args := m.Called(noteId)

	return args.Error(0)
}
func (m *MockedReferences) ImportMany(refs []ReferencesDocument) (int, error) {
	args := m.Called(refs)

	return args.Int(0), args.Error(1)
}

func (m *MockedReferences) Export(noteIds []string) ([]ReferencesDocument, error) {
	args := m.Called()

	return args.Get(0).([]ReferencesDocument), args.Error(1)
}
