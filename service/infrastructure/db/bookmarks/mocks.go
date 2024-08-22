package bookmarks

import "github.com/stretchr/testify/mock"

type MockedBookmarks struct {
	mock.Mock
}

func (m *MockedBookmarks) SearchTags(userId, searchKeyword string, limit int) ([]TagsResult, error) {
	args := m.Called(userId, searchKeyword, limit)

	return args.Get(0).([]TagsResult), args.Error(1)
}
func (m *MockedBookmarks) Count(userId string) (int64, error) {
	args := m.Called(userId)

	return args.Get(0).(int64), args.Error(1)
}
func (m *MockedBookmarks) List(userId, searchKeyword string, page, pageSize int, sortBy, sortDirection string, tags []string) ([]BookmarkDocument, int, error) {
	args := m.Called(userId, searchKeyword, page, pageSize, sortBy, sortDirection, tags)

	return args.Get(0).([]BookmarkDocument), args.Int(1), args.Error(2)
}
func (m *MockedBookmarks) GetOne(Id string) (BookmarkDocument, error) {
	args := m.Called(Id)

	return args.Get(0).(BookmarkDocument), args.Error(1)
}
func (m *MockedBookmarks) GetBookmarks(ids, fields []string) ([]BookmarkDocument, error) {
	args := m.Called(ids, fields)

	return args.Get(0).([]BookmarkDocument), args.Error(1)
}
func (m *MockedBookmarks) InsertOne(userId, URL, title string, tags []string) (BookmarkDocument, error) {
	args := m.Called(userId, URL, title, tags)

	return args.Get(0).(BookmarkDocument), args.Error(1)
}
func (m *MockedBookmarks) UpdateOne(userId, Id string, updates interface{}) error {
	args := m.Called(userId, Id, updates)

	return args.Error(0)
}
func (m *MockedBookmarks) DeleteOne(Id string) error {
	args := m.Called(Id)

	return args.Error(0)
}
