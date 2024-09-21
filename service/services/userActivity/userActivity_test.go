package userActivityService

import (
	"fmt"
	"testing"

	"github.com/alperklc/the-zula/service/infrastructure/cache"
	useractivity "github.com/alperklc/the-zula/service/infrastructure/db/userActivity"
	bookmarksService "github.com/alperklc/the-zula/service/services/bookmarks"
	notesService "github.com/alperklc/the-zula/service/services/notes"
	usersService "github.com/alperklc/the-zula/service/services/users"
	"github.com/stretchr/testify/assert"
)

func TestUserActivity(t *testing.T) {
	t.Run("it wont create user activity if the resource type is missing", func(t *testing.T) {
		// arrange mocks
		mostVisitedCache := new(cache.MockedCache[[]useractivity.UsageStatisticsEntry])
		activityGraphCache := new(cache.MockedCache[[]useractivity.ActivityGraphEntry])

		userRepository := new(usersService.MockedUser)
		userRepository.On("GetUser", "userId").Return(usersService.User{ID: "userID"}, nil)
		activityController := NewService(activityGraphCache, mostVisitedCache, userRepository, nil, nil, nil)

		// act
		_, error := activityController.Create("userId", "", "CREATE", "NOTE_ID")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("MISSING_FIELDS"))
	})

	t.Run("it wont create user activity if the action is missing", func(t *testing.T) {
		// arrange mocks
		mostVisitedCache := new(cache.MockedCache[[]useractivity.UsageStatisticsEntry])
		activityGraphCache := new(cache.MockedCache[[]useractivity.ActivityGraphEntry])

		userRepository := new(usersService.MockedUser)
		userRepository.On("GetUser", "userId").Return(usersService.User{ID: "userID"}, nil)
		activityController := NewService(activityGraphCache, mostVisitedCache, userRepository, nil, nil, nil)

		// act
		_, error := activityController.Create("userId", "NOTE", "", "NOTE_ID")

		// assert
		assert := assert.New(t)
		assert.Equal(error, fmt.Errorf("MISSING_FIELDS"))
	})

	t.Run("it calls insertOne method of the database while creating user activity", func(t *testing.T) {
		// arrange mocks
		mostVisitedCache := new(cache.MockedCache[[]useractivity.UsageStatisticsEntry])
		activityGraphCache := new(cache.MockedCache[[]useractivity.ActivityGraphEntry])

		userRepository := new(usersService.MockedUser)
		userRepository.On("GetUser", "userId").Return(usersService.User{ID: "userID"}, nil)
		userActivityRepository := new(useractivity.MockedUserActivity)
		userActivityRepository.On("InsertOne", "userId", "NOTE", "CREATE", "NOTE_ID").Return(useractivity.UserActivityDocument{}, nil)

		activityController := NewService(activityGraphCache, mostVisitedCache, userRepository, userActivityRepository, nil, nil)

		// act
		_, error := activityController.Create("userId", "NOTE", "CREATE", "NOTE_ID")

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
	})

	t.Run("it calls list method of the database while listing user activity", func(t *testing.T) {
		// arrange mocks
		mostVisitedCache := new(cache.MockedCache[[]useractivity.UsageStatisticsEntry])
		activityGraphCache := new(cache.MockedCache[[]useractivity.ActivityGraphEntry])

		userRepository := new(usersService.MockedUser)
		userRepository.On("GetUser", "userId").Return(usersService.User{ID: "userId"}, nil)
		userActivityRepository := new(useractivity.MockedUserActivity)
		userActivityRepository.On("List", "userId", 1, 10, "timestamp", "desc").Return(useractivity.UserActivityPage{}, nil)

		activityController := NewService(activityGraphCache, mostVisitedCache, userRepository, userActivityRepository, nil, nil)

		// act
		_, error := activityController.List("userId", 1, 10, "timestamp", "desc")

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
	})

	t.Run("it fetches most visited content from cache", func(t *testing.T) {
		// arrange mocks
		mostVisitedCache := new(cache.MockedCache[[]useractivity.UsageStatisticsEntry])
		activityGraphCache := new(cache.MockedCache[[]useractivity.ActivityGraphEntry])

		userRepository := new(usersService.MockedUser)
		userRepository.On("GetUser", "user").Return(usersService.User{ID: "user"}, nil)
		mostVisitedContent := []useractivity.UsageStatisticsEntry{{ObjectID: "asd"}, {ObjectID: "asdasd"}}
		mostVisitedCache.On("Read", CACHE_PREFIX_MOST_VISITED+"user").Return(&mostVisitedContent)

		activityController := NewService(activityGraphCache, mostVisitedCache, userRepository, nil, nil, nil)

		// act
		mostVisited, error := activityController.GetMostVisited("user")

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
		assert.Len(mostVisited, len(mostVisitedContent))
	})

	t.Run("it queries for insights for dashboard", func(t *testing.T) {
		// arrange mocks
		mostVisitedCache := new(cache.MockedCache[[]useractivity.UsageStatisticsEntry])
		activityGraphCache := new(cache.MockedCache[[]useractivity.ActivityGraphEntry])

		bs := new(bookmarksService.MockedBookmarkService)
		ns := new(notesService.MockedNoteService)
		ns.On("GetStatistics", "user").Return(notesService.Statistics{Count: 5}, nil)
		bs.On("GetStatistics", "user").Return(bookmarksService.Statistics{Count: 3}, nil)

		userRepository := new(usersService.MockedUser)
		userRepository.On("GetUser", "user").Return(usersService.User{ID: "user"}, nil)
		mostVisitedContent := []useractivity.UsageStatisticsEntry{{ObjectID: "asd"}, {ObjectID: "asdasd"}}
		lastVisitedContent := []useractivity.UsageStatisticsEntry{{ObjectID: "note1"}, {ObjectID: "note2"}}
		activityGraphContent := []useractivity.ActivityGraphEntry{{Date: "2021-02-01"}}

		mostVisitedCache.On("Read", CACHE_PREFIX_MOST_VISITED+"user").Return(&mostVisitedContent)
		activityGraphCache.On("Read", CACHE_PREFIX_ACTIVITY_GRAPH+"user").Return(&activityGraphContent)
		userActivityRepository := new(useractivity.MockedUserActivity)
		userActivityRepository.On("GetLastVisitedContent", "user", -7, []string{}).Return(lastVisitedContent, nil)
		userActivityRepository.On("GetIdsOfDeletedEntries", "user", -7).Return([]string{}, nil)

		activityController := NewService(activityGraphCache, mostVisitedCache, userRepository, userActivityRepository, ns, bs)

		// act
		activityGraph, mostVisited, lastVisited, nrOfNotes, nrOfBookmarks, error := activityController.GetInsightsForDashboard("user")

		// assert
		assert := assert.New(t)
		assert.Equal(error, nil)
		assert.Len(mostVisited, len(mostVisitedContent))
		assert.Len(lastVisited, len(lastVisitedContent))
		assert.Len(activityGraph, len(activityGraphContent))
		assert.Greater(nrOfNotes, int64(0))
		assert.Greater(nrOfBookmarks, int64(0))
	})
}
