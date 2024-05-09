package useractivity

import "time"

type UserActivityDocument struct {
	UserUID      string    `json:"userUID" bson:"userUID"`
	ResourceType string    `json:"resourceType" bson:"resourceType"`
	Action       string    `json:"action" bson:"action"`
	ObjectUID    string    `json:"objectUID" bson:"objectUID"`
	Timestamp    time.Time `json:"timestamp" bson:"timestamp"`
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
	ObjectUID    string    `json:"objectUID" bson:"objectUID"`
	Timestamp    time.Time `json:"timestamp" bson:"timestamp"`
	Count        int       `json:"count" bson:"count"`
}

type ActivityGraphEntry struct {
	Date  string `json:"date" bson:"_id"`
	Count int    `json:"count" bson:"count"`
}
