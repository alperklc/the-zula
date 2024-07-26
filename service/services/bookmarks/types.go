package bookmarksService

import (
	"time"

	"github.com/alperklc/the-zula/service/infrastructure/db/pageContent"
)

type Bookmark struct {
	ShortId     string      `json:"shortId"`
	URL         string      `json:"url"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	UpdatedBy   string      `json:"updatedBy"`
	CreatedBy   string      `json:"createdBy"`
	CreatedAt   time.Time   `json:"createdAt"`
	Title       string      `json:"title"`
	FaviconURL  string      `json:"faviconUrl"`
	Tags        []string    `json:"tags"`
	PageContent PageContent `json:"pageContent"`
}

func (b *Bookmark) AddPageContent(input pageContent.PageContentDocument) Bookmark {
	pc := PageContent{
		ID:        input.Id,
		URL:       input.URL,
		Title:     input.Title,
		Author:    input.Author,
		Length:    input.Length,
		SiteName:  input.SiteName,
		Image:     input.Image,
		Favicon:   input.Favicon,
		MDContent: input.MDContent,
	}
	return Bookmark{
		ShortId:     b.ShortId,
		URL:         b.URL,
		UpdatedAt:   b.UpdatedAt,
		UpdatedBy:   b.UpdatedBy,
		CreatedBy:   b.CreatedBy,
		CreatedAt:   b.CreatedAt,
		Title:       b.Title,
		FaviconURL:  b.FaviconURL,
		PageContent: pc,
		Tags:        b.Tags,
	}
}

type PaginationMeta struct {
	Count         int    `json:"count"`
	Query         string `json:"query"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
	SortBy        string `json:"sortBy"`
	SortDirection string `json:"sortDirection"`
	Range         string `json:"range"`
}

type BookmarksPage struct {
	Meta  PaginationMeta `json:"meta"`
	Items []Bookmark     `json:"items"`
}

type Statistics struct {
	Count int64 `json:"count"`
}

type PageContent struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Length      int    `json:"int"`
	Excerpt     string `json:"excerpt"`
	SiteName    string `json:"siteName"`
	Image       string `json:"image"`
	Favicon     string `json:"favicon"`
	TextContent string `json:"textContent"`
	HTMLContent string `json:"htmlContent"`
	MDContent   string `json:"mdContent"`
}

func (p *PageContent) ConvertToDBPageContent() pageContent.PageContent {
	return pageContent.PageContent{
		URL:         p.URL,
		Title:       p.Title,
		Author:      p.Author,
		Length:      p.Length,
		Excerpt:     p.Excerpt,
		SiteName:    p.SiteName,
		Image:       p.Image,
		Favicon:     p.Favicon,
		MDContent:   p.MDContent,
		HTMLContent: p.HTMLContent,
	}
}
