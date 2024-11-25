package userActivityService

import (
	"fmt"

	"github.com/alperklc/the-zula/service/infrastructure/cache"
	useractivity "github.com/alperklc/the-zula/service/infrastructure/db/userActivity"
	bookmarksService "github.com/alperklc/the-zula/service/services/bookmarks"
	notesService "github.com/alperklc/the-zula/service/services/notes"
	usersService "github.com/alperklc/the-zula/service/services/users"
	"golang.org/x/exp/maps"
)

const (
	CACHE_PREFIX_MOST_VISITED   = "MOST_VISITED"
	CACHE_PREFIX_ACTIVITY_GRAPH = "ACTIVITY_GRAPH"
)

type UserActivityService interface {
	Create(userID, resourceType, action, objectID string) (UserActivityDocument, error)
	List(userID string, page, pageSize int, sortBy, sortDirection string) (UserActivityPage, error)
	GetMostVisited(userID string) ([]useractivity.UsageStatisticsEntry, error)
	GetLastVisited(userID string) ([]useractivity.UsageStatisticsEntry, error)
	GetInsightsForDashboard(userID string) ([]useractivity.ActivityGraphEntry, []useractivity.UsageStatisticsEntry, []useractivity.UsageStatisticsEntry, int64, int64, error)
}

type datasources struct {
	activityGraphCache cache.CacheInterface[[]useractivity.ActivityGraphEntry]
	mostVisitedCache   cache.CacheInterface[[]useractivity.UsageStatisticsEntry]
	user               usersService.UsersService
	notes              notesService.NoteService
	bookmarks          bookmarksService.BookmarkService
	useractivity       useractivity.Collection
}

func NewService(agc cache.CacheInterface[[]useractivity.ActivityGraphEntry], mvc cache.CacheInterface[[]useractivity.UsageStatisticsEntry], u usersService.UsersService, ua useractivity.Collection, n notesService.NoteService, b bookmarksService.BookmarkService) UserActivityService {
	return &datasources{
		activityGraphCache: agc,
		mostVisitedCache:   mvc,
		user:               u,
		useractivity:       ua,
		notes:              n,
		bookmarks:          b,
	}
}

func (d *datasources) Create(userID, resourceType, action, objectID string) (UserActivityDocument, error) {
	foundUser, getUserError := d.user.GetUser(userID)
	if foundUser.ID == "" || getUserError != nil {
		return UserActivityDocument{}, fmt.Errorf("USER_NOT_FOUND")
	}

	if resourceType == "" || action == "" {
		return UserActivityDocument{}, fmt.Errorf("MISSING_FIELDS")
	}

	userActivity, createErr := d.useractivity.InsertOne(
		userID,
		resourceType,
		action,
		objectID,
	)

	return UserActivityDocument(userActivity), createErr
}

func (d *datasources) List(userID string, page, pageSize int, sortBy, sortDirection string) (UserActivityPage, error) {
	foundUser, getUserError := d.user.GetUser(userID)
	if foundUser.ID == "" || getUserError != nil {
		return UserActivityPage{}, fmt.Errorf("USER_NOT_FOUND")
	}

	list, listErr := d.useractivity.List(userID, page, pageSize, sortBy, sortDirection)

	var items []UserActivityDocument = make([]UserActivityDocument, 0, len(list.Items))
	for _, b := range list.Items {
		note := UserActivityDocument{
			UserID:       b.UserID,
			ObjectID:     b.ObjectID,
			Action:       b.Action,
			ResourceType: b.ResourceType,
			Timestamp:    b.Timestamp,
		}
		items = append(items, note)
	}

	return UserActivityPage{Meta: PaginationMeta(list.Meta), Items: items}, listErr
}

func (d *datasources) GetMostVisited(userID string) ([]useractivity.UsageStatisticsEntry, error) {
	obj := d.mostVisitedCache.Read(CACHE_PREFIX_MOST_VISITED + userID)
	if obj != nil {
		return *obj, nil
	}

	uidsOfDeletedItems, _ := d.useractivity.GetIdsOfDeletedEntries(userID, -7)
	mostVisitedItems, errMostVisitedItems := d.useractivity.GetMostVisitedContent(userID, -7, uidsOfDeletedItems)
	if errMostVisitedItems != nil {
		return nil, errMostVisitedItems
	}

	var idsMap map[string]string = make(map[string]string)
	for _, item := range mostVisitedItems {
		idsMap[item.ObjectID] = ""
	}

	bookmarksFound, _ := d.bookmarks.GetBookmarks(maps.Keys(idsMap), []string{"title", "shortId"})
	notesFound, _ := d.notes.GetNotes(maps.Keys(idsMap), []string{"title", "shortId"})

	for i := range mostVisitedItems {
		if mostVisitedItems[i].ResourceType == "BOOKMARK" {
			mostVisitedItems[i].Title = bookmarksFound[mostVisitedItems[i].ObjectID].Title
		} else if mostVisitedItems[i].ResourceType == "NOTE" {
			mostVisitedItems[i].Title = notesFound[mostVisitedItems[i].ObjectID].Title
		}
	}

	d.mostVisitedCache.Write(CACHE_PREFIX_MOST_VISITED+userID, mostVisitedItems)
	return mostVisitedItems, nil
}

func (d *datasources) GetLastVisited(userID string) ([]useractivity.UsageStatisticsEntry, error) {
	uidsOfDeletedItems, _ := d.useractivity.GetIdsOfDeletedEntries(userID, -7)

	lastVisitedItems, errLastVisitedItems := d.useractivity.GetLastVisitedContent(userID, -7, uidsOfDeletedItems)
	if errLastVisitedItems != nil {
		return nil, errLastVisitedItems
	}

	return lastVisitedItems, nil
}

func (d *datasources) GroupActivitiesByDate(userID string) ([]useractivity.ActivityGraphEntry, error) {
	obj := d.activityGraphCache.Read(CACHE_PREFIX_ACTIVITY_GRAPH + userID)
	if obj != nil {
		return *obj, nil
	}
	activityGraph, errActivities := d.useractivity.GroupActivitiesByDate(userID)
	if errActivities != nil {
		return nil, errActivities
	}

	d.activityGraphCache.Write(CACHE_PREFIX_ACTIVITY_GRAPH+userID, activityGraph)
	return activityGraph, nil
}

func (d *datasources) GetInsightsForDashboard(userID string) ([]useractivity.ActivityGraphEntry, []useractivity.UsageStatisticsEntry, []useractivity.UsageStatisticsEntry, int64, int64, error) {
	foundUser, getUserError := d.user.GetUser(userID)
	if foundUser.ID == "" || getUserError != nil {
		return nil, nil, nil, 0, 0, fmt.Errorf("USER_NOT_FOUND")
	}
	activityGraph, errActivities := d.GroupActivitiesByDate(userID)
	if errActivities != nil {
		return nil, nil, nil, 0, 0, errActivities
	}

	mostVisited, errMostVisited := d.GetMostVisited(userID)
	if errMostVisited != nil {
		return nil, nil, nil, 0, 0, errMostVisited
	}
	fmt.Println(mostVisited)

	lastVisited, errLastVisited := d.GetLastVisited(userID)
	if errLastVisited != nil {
		return nil, nil, nil, 0, 0, errLastVisited
	}

	nrOfNotes, errNotesStats := d.notes.GetStatistics(userID)
	if errNotesStats != nil {
		return nil, nil, nil, 0, 0, errNotesStats
	}

	nrOfBookmarks, errBookmarksStats := d.bookmarks.GetStatistics(userID)
	if errBookmarksStats != nil {
		return nil, nil, nil, 0, 0, errBookmarksStats
	}

	return activityGraph, mostVisited, lastVisited, nrOfNotes.Count, nrOfBookmarks.Count, nil
}
