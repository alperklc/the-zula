package userActivityService

import "time"

type UserActivityDocument struct {
	Id           string    `bson:"_id,omitempty"`
	UserID       string    `json:"userID"`
	ResourceType string    `json:"resourceType"`
	Action       string    `json:"action"`
	ObjectID     string    `json:"objectID"`
	Timestamp    time.Time `json:"timestamp"`
}

type PaginationMeta struct {
	Count         int    `json:"count"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
	SortBy        string `json:"sortBy"`
	SortDirection string `json:"sortDirection"`
	Range         string `json:"range"`
}

type UserActivityPage struct {
	Meta  PaginationMeta         `json:"meta"`
	Items []UserActivityDocument `json:"items"`
}

type UsageStatisticsEntry struct {
	ResourceType string    `json:"resourceType" bson:"resourceType"`
	ObjectID     string    `json:"objectID" bson:"objectID"`
	Timestamp    time.Time `json:"timestamp" bson:"timestamp"`
	Count        int       `json:"count" bson:"count"`
}

type UserActivityInsights struct {
	MostReadContent    UsageStatisticsEntry `json:"mostReadContent"`
	LastCreatedContent UsageStatisticsEntry `json:"lastCreatedContent"`
	LastUpdatedContent UsageStatisticsEntry `json:"lastUpdatedContent"`
}
