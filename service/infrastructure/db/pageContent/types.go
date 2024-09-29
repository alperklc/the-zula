package pageContent

import (
	"time"
)

type PageContentDocument struct {
	Id        string    `bson:"_id,omitempty"`
	CreatedAt time.Time `bson:"createdAt"`
	URL       string    `bson:"url"`
	Title     string    `bson:"title"`
	Author    string    `bson:"author"`
	Length    int       `bson:"length"`
	SiteName  string    `bson:"siteName"`
	Image     string    `bson:"image"`
	Favicon   string    `bson:"favicon"`
	MDContent string    `bson:"mdContent"`
}

type PageContent struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Length      int    `json:"int"`
	Excerpt     string `json:"excerpt"`
	SiteName    string `json:"siteName"`
	Image       string `json:"image"`
	Favicon     string `json:"favicon"`
	HTMLContent string `json:"htmlContent"`
	MDContent   string `json:"mdContent"`
}

func (n PageContentDocument) GetId() string {
	return n.Id
}
