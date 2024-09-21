package notesService

import (
	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/stretchr/testify/mock"
)

type MockedNoteService struct {
	mock.Mock
}

func (m *MockedNoteService) SearchTags(userId, searchKeyword string, limit int) ([]notes.TagsResult, error) {
	args := m.Called(userId, searchKeyword, limit)
	return args.Get(0).([]notes.TagsResult), args.Error(1)
}

func (m *MockedNoteService) GetStatistics(userId string) (Statistics, error) {
	args := m.Called(userId)
	return args.Get(0).(Statistics), args.Error(1)
}

func (m *MockedNoteService) ListNotes(userId string, searchKeyword *string, page, pageSize *int, sortBy, sortDirection *string, tags *[]string) (NotesPage, error) {
	args := m.Called(userId, searchKeyword, page, pageSize, sortBy, sortDirection, tags)
	return args.Get(0).(NotesPage), args.Error(1)
}

func (m *MockedNoteService) CreateNote(userId, clientId string, title, content *string, tags *[]string) (Note, error) {
	args := m.Called(userId, clientId, title, content, tags)
	return args.Get(0).(Note), args.Error(1)
}

func (m *MockedNoteService) UpdateNote(noteId, userId, clientId string, update map[string]interface{}) error {
	args := m.Called(noteId, userId, clientId, update)
	return args.Error(0)
}

func (m *MockedNoteService) GetNote(noteId, userId, clientId string, params GetNoteParams) (Note, error) {
	args := m.Called(noteId, userId, clientId, params)
	return args.Get(0).(Note), args.Error(1)
}

func (m *MockedNoteService) GetNotes(noteIds, fields []string) (map[string]Note, error) {
	args := m.Called(noteIds, fields)
	return args.Get(0).(map[string]Note), args.Error(1)
}

func (m *MockedNoteService) DeleteNote(noteId, userId, clientId string) error {
	args := m.Called(noteId, userId, clientId)
	return args.Error(0)
}

func (m *MockedNoteService) GetDraftOfNote(userId, noteId string) (Note, error) {
	args := m.Called(userId, noteId)
	return args.Get(0).(Note), args.Error(1)
}

func (m *MockedNoteService) UpdateDraft(userId, noteId, title, content string, tags []string) error {
	args := m.Called(userId, noteId, title, content, tags)
	return args.Error(0)
}

func (m *MockedNoteService) DeleteDraft(userId, noteId string) error {
	args := m.Called(userId, noteId)
	return args.Error(0)
}

func (m *MockedNoteService) ListNotesChanges(userId, noteId string, page, pageSize *int) (NotesChangesPage, error) {
	args := m.Called(userId, noteId, page, pageSize)
	return args.Get(0).(NotesChangesPage), args.Error(1)
}

func (m *MockedNoteService) GetNotesChange(noteId, shortId string) (NotesChanges, error) {
	args := m.Called(noteId, shortId)
	return args.Get(0).(NotesChanges), args.Error(1)
}
