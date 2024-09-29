package bookmarks

import "time"

type BookmarkDocument struct {
	Id         string    `bson:"_id,omitempty"`
	ShortId    string    `bson:"shortId"`
	URL        string    `bson:"url"`
	UpdatedAt  time.Time `bson:"updatedAt"`
	UpdatedBy  string    `bson:"updatedBy"`
	CreatedBy  string    `bson:"createdBy"`
	CreatedAt  time.Time `bson:"createdAt"`
	Title      string    `bson:"title"`
	FaviconURL string    `bson:"faviconUrl"`
	Tags       []string  `bson:"tags"`
}

func (n BookmarkDocument) GetId() string {
	return n.Id
}

type BookmarkFacets struct {
	Count int `json:"count"`
}

type BookmarksAggregation struct {
	Facets BookmarkFacets     `json:"facets"`
	Items  []BookmarkDocument `json:"items"`
}

type TagsResult struct {
	Value string `bson:"_id" json:"value,omitempty"`
	Count int    `bson:"count" json:"count"`
}
