package notes

import (
	"github.com/stretchr/testify/mock"
)

type MockedNotes struct {
	mock.Mock
}

func (m *MockedNotes) SearchTags(user, searchKeyword string, limit int) ([]TagsResult, error) {
	args := m.Called(user, searchKeyword, limit)

	return args.Get(0).([]TagsResult), args.Error(1)
}
func (m *MockedNotes) Count(user string) (int64, error) {
	args := m.Called(user)

	return args.Get(0).(int64), args.Error(1)
}
func (m *MockedNotes) List(user, searchKeyword string, page, pageSize int, sortBy, sortDirection string, tags []string) ([]NoteDocument, int, error) {
	args := m.Called(user, searchKeyword, page, pageSize, sortBy, sortDirection, tags)

	return args.Get(0).([]NoteDocument), args.Int(1), args.Error(2)
}
func (m *MockedNotes) GetNotes(ids, fields []string) ([]NoteDocument, error) {
	args := m.Called(ids, fields)

	return args.Get(0).([]NoteDocument), args.Error(1)
}
func (m *MockedNotes) GetOne(id string) (NoteDocument, error) {
	args := m.Called(id)

	return args.Get(0).(NoteDocument), args.Error(1)
}
func (m *MockedNotes) InsertOne(user, title, content string, tags []string) (NoteDocument, error) {
	args := m.Called(user, title, content, tags)

	return args.Get(0).(NoteDocument), args.Error(1)
}
func (m *MockedNotes) UpdateOne(user, id string, updates interface{}) error {
	args := m.Called(user, id, updates)

	return args.Error(0)
}
func (m *MockedNotes) DeleteOne(id string) error {
	args := m.Called(id)

	return args.Error(0)
}
func (m *MockedNotes) ImportMany(notes []NoteDocument) (int, error) {
	args := m.Called(notes)

	return args.Int(0), args.Error(1)
}
func (m *MockedNotes) ExportForUser(userId string) ([]NoteDocument, error) {
	args := m.Called(userId)

	return args.Get(0).([]NoteDocument), args.Error(1)
}
