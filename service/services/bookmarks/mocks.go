package bookmarksService

import (
	"github.com/alperklc/the-zula/service/infrastructure/db/bookmarks"
	"github.com/stretchr/testify/mock"
)

type MockedBookmarkService struct {
	mock.Mock
}

func (m *MockedBookmarkService) SearchTags(userID, searchKeyword string, limit int) ([]bookmarks.TagsResult, error) {
	args := m.Called(userID, searchKeyword, limit)
	return args.Get(0).([]bookmarks.TagsResult), args.Error(1)
}

func (m *MockedBookmarkService) GetStatistics(userID string) (Statistics, error) {
	args := m.Called(userID)
	return args.Get(0).(Statistics), args.Error(1)
}

func (m *MockedBookmarkService) ListBookmarks(userID string, searchKeyword *string, page, pageSize *int, sortBy, sortDirection *string, tags *[]string) (BookmarksPage, error) {
	args := m.Called(userID, searchKeyword, page, pageSize, sortBy, sortDirection, tags)
	return args.Get(0).(BookmarksPage), args.Error(1)
}

func (m *MockedBookmarkService) CreateBookmark(userID, clientId, URL, title string, tags *[]string) (Bookmark, error) {
	args := m.Called(userID, clientId, URL, title, tags)
	return args.Get(0).(Bookmark), args.Error(1)
}

func (m *MockedBookmarkService) UpdateBookmark(bookmarkID, userID, clientId string, update map[string]interface{}) error {
	args := m.Called(bookmarkID, userID, clientId, update)
	return args.Error(0)
}

func (m *MockedBookmarkService) GetBookmark(bookmarkID, userID, clientId string) (Bookmark, error) {
	args := m.Called(bookmarkID, userID, clientId)
	return args.Get(0).(Bookmark), args.Error(1)
}

func (m *MockedBookmarkService) GetBookmarks(ids, fields []string) (map[string]Bookmark, error) {
	args := m.Called(ids, fields)
	return args.Get(0).(map[string]Bookmark), args.Error(1)
}

func (m *MockedBookmarkService) DeleteBookmark(bookmarkID, userID, clientId string) error {
	args := m.Called(bookmarkID, userID, clientId)
	return args.Error(0)
}

func (m *MockedBookmarkService) GetPageContentOfBookmark(bookmarkID string) (PageContent, error) {
	args := m.Called(bookmarkID)
	return args.Get(0).(PageContent), args.Error(1)
}

func (m *MockedBookmarkService) ParsePageContentOfBookmark(bookmark bookmarks.BookmarkDocument) (PageContent, error) {
	args := m.Called(bookmark)
	return args.Get(0).(PageContent), args.Error(1)
}
