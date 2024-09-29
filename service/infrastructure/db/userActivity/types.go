package useractivity

import "time"

type UserActivityDocument struct {
	Id           string    `bson:"_id,omitempty"`
	UserID       string    `json:"userID" bson:"userID"`
	ResourceType string    `json:"resourceType" bson:"resourceType"`
	Action       string    `json:"action" bson:"action"`
	ObjectID     string    `json:"objectID" bson:"objectID"`
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
	ObjectID     string    `json:"objectID" bson:"objectID"`
	Timestamp    time.Time `json:"timestamp" bson:"timestamp"`
	Count        int       `json:"count" bson:"count"`
	Title        string    `json:"title" bson:"title"`
}

type ActivityGraphEntry struct {
	Date  string `json:"date" bson:"_id"`
	Count int    `json:"count" bson:"count"`
}

func (n UserActivityDocument) GetId() string {
	return n.Id
}
