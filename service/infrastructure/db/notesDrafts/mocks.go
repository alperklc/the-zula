package notesDrafts

import (
	"github.com/stretchr/testify/mock"
)

type MockedNotesDrafts struct {
	mock.Mock
}

func (m *MockedNotesDrafts) CheckExistence(ids []string) (map[string]bool, error) {
	args := m.Called(ids)

	return args.Get(0).(map[string]bool), args.Error(1)
}

func (m *MockedNotesDrafts) GetOne(id string) (NoteDraftDocument, error) {
	args := m.Called(id)

	return args.Get(0).(NoteDraftDocument), args.Error(1)
}

func (m *MockedNotesDrafts) UpsertOne(id, title, content string, tags []string) error {
	args := m.Called(id, title, content, tags)

	return args.Error(0)
}

func (m *MockedNotesDrafts) DeleteOne(id string) error {
	args := m.Called(id)

	return args.Error(0)
}
