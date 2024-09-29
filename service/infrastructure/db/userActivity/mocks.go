package useractivity

import (
	"github.com/stretchr/testify/mock"
)

type MockedUserActivity struct {
	mock.Mock
}

func (m *MockedUserActivity) List(userID string, page, pageSize int, sortBy, sortDirection string) (UserActivityPage, error) {
	args := m.Called(userID, page, pageSize, sortBy, sortDirection)
	return args.Get(0).(UserActivityPage), args.Error(1)
}

func (m *MockedUserActivity) InsertOne(userID, resourceType, action, objectID string) (UserActivityDocument, error) {
	args := m.Called(userID, resourceType, action, objectID)
	return args.Get(0).(UserActivityDocument), args.Error(1)
}

func (m *MockedUserActivity) GetMostVisitedContent(userID string, daysSince int, IDsOfExcludedEntries []string) ([]UsageStatisticsEntry, error) {
	args := m.Called(userID, daysSince, IDsOfExcludedEntries)
	return args.Get(0).([]UsageStatisticsEntry), args.Error(1)
}

func (m *MockedUserActivity) GetLastVisitedContent(userID string, daysSince int, IDsOfExcludedEntries []string) ([]UsageStatisticsEntry, error) {
	args := m.Called(userID, daysSince, IDsOfExcludedEntries)
	return args.Get(0).([]UsageStatisticsEntry), args.Error(1)
}

func (m *MockedUserActivity) GetIdsOfDeletedEntries(userID string, daysAgo int) ([]string, error) {
	args := m.Called(userID, daysAgo)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockedUserActivity) GroupActivitiesByDate(userID string) ([]ActivityGraphEntry, error) {
	args := m.Called(userID)
	return args.Get(0).([]ActivityGraphEntry), args.Error(1)
}
func (m *MockedUserActivity) ImportMany(refs []UserActivityDocument) (int, error) {
	args := m.Called(refs)

	return args.Int(0), args.Error(1)
}
func (m *MockedUserActivity) ExportForUser(userId string) ([]UserActivityDocument, error) {
	args := m.Called(userId)
	return args.Get(0).([]UserActivityDocument), args.Error(1)
}
