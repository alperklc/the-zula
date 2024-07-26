package bookmarks

import "github.com/stretchr/testify/mock"

type mockedBookmarks struct {
	mock.Mock
}

func (m *mockedBookmarks) SearchTags(userUID, searchKeyword string, limit int) ([]TagsResult, error) {
	args := m.Called(userUID, searchKeyword, limit)

	return args.Get(0).([]TagsResult), args.Error(1)
}
func (m *mockedBookmarks) Count(userUID string) (int64, error) {
	args := m.Called(userUID)

	return args.Get(0).(int64), args.Error(1)
}
func (m *mockedBookmarks) List(userUID, searchKeyword string, page, pageSize int, sortBy, sortDirection string, tags []string) ([]BookmarkDocument, int, error) {
	args := m.Called(userUID, searchKeyword, page, pageSize, sortBy, sortDirection, tags)

	return args.Get(0).([]BookmarkDocument), args.Int(1), args.Error(2)
}
func (m *mockedBookmarks) GetOne(UID string) (BookmarkDocument, error) {
	args := m.Called(UID)

	return args.Get(0).(BookmarkDocument), args.Error(1)
}
func (m *mockedBookmarks) InsertOne(userId, URL, title string, tags []string) (BookmarkDocument, error) {
	args := m.Called(userId, URL, title, tags)

	return args.Get(0).(BookmarkDocument), args.Error(1)
}
func (m *mockedBookmarks) UpdateOne(userUID, UID string, updates interface{}) error {
	args := m.Called(userUID, UID, updates)

	return args.Error(0)
}
func (m *mockedBookmarks) DeleteOne(UID string) error {
	args := m.Called(UID)

	return args.Error(0)
}
