package useractivity

import (
	"github.com/stretchr/testify/mock"
)

type MockedUserActivity struct {
	mock.Mock
}

func (m *MockedUserActivity) List(userUID string, page, pageSize int, sortBy, sortDirection string) (UserActivityPage, error) {
	args := m.Called(userUID, page, pageSize, sortBy, sortDirection)
	return args.Get(0).(UserActivityPage), args.Error(1)
}

func (m *MockedUserActivity) InsertOne(userUID, resourceType, action, objectUID string) (UserActivityDocument, error) {
	args := m.Called(userUID, resourceType, action, objectUID)
	return args.Get(0).(UserActivityDocument), args.Error(1)
}

func (m *MockedUserActivity) GetMostVisitedContent(userUID string, daysSince int, uidsOfExcludedEntries []string) ([]UsageStatisticsEntry, error) {
	args := m.Called(userUID, daysSince, uidsOfExcludedEntries)
	return args.Get(0).([]UsageStatisticsEntry), args.Error(1)
}

func (m *MockedUserActivity) GetLastVisitedContent(userUID string, daysSince int, uidsOfExcludedEntries []string) ([]UsageStatisticsEntry, error) {
	args := m.Called(userUID, daysSince, uidsOfExcludedEntries)
	return args.Get(0).([]UsageStatisticsEntry), args.Error(1)
}

func (m *MockedUserActivity) GetUidsOfDeletedEntries(userUID string, daysAgo int) ([]string, error) {
	args := m.Called(userUID, daysAgo)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockedUserActivity) GroupActivitiesByDate(userUID string) ([]ActivityGraphEntry, error) {
	args := m.Called(userUID)
	return args.Get(0).([]ActivityGraphEntry), args.Error(1)
}
