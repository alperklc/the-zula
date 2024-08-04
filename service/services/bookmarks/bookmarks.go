package bookmarksService

import (
	"fmt"
	"math"

	"github.com/alperklc/the-zula/service/infrastructure/db/bookmarks"
	"github.com/alperklc/the-zula/service/infrastructure/db/pageContent"
	mqpublisher "github.com/alperklc/the-zula/service/infrastructure/messageQueue/publisher"
	"github.com/alperklc/the-zula/service/infrastructure/webScraper"
	usersService "github.com/alperklc/the-zula/service/services/users"
	"github.com/alperklc/the-zula/service/utils"
	gonanoid "github.com/matoous/go-nanoid"
)

type BookmarkService interface {
	SearchTags(userID, searchKeyword string, limit int) ([]bookmarks.TagsResult, error)
	GetStatistics(userID string) (Statistics, error)
	ListBookmarks(userID string, searchKeyword *string, page, pageSize *int, sortBy, sortDirection *string, tags *[]string) (BookmarksPage, error)
	CreateBookmark(userID, clientId, URL, title string, tags *[]string) (Bookmark, error)
	UpdateBookmark(bookmarkID, userID, clientId string, update map[string]interface{}) error
	GetBookmark(bookmarkID, userID, clientId string) (Bookmark, error)
	DeleteBookmark(bookmarkID, userID, clientId string) error
	GetPageContentOfBookmark(bookmarkID string) (PageContent, error)
	ParsePageContentOfBookmark(bookmark bookmarks.BookmarkDocument) (PageContent, error)
}

type datasources struct {
	users       usersService.UsersService
	bookmarks   bookmarks.Collection
	pageContent pageContent.Collection
	webScraper  webScraper.WebScraper
	mqpublisher mqpublisher.MessagePublisher
}

func NewService(u usersService.UsersService, b bookmarks.Collection, pc pageContent.Collection, w webScraper.WebScraper, mqp mqpublisher.MessagePublisher) BookmarkService {
	return &datasources{users: u, bookmarks: b, pageContent: pc, webScraper: w, mqpublisher: mqp}
}

func getPaginationRange(count, page, pageSize int) string {
	if count == 0 {
		return ""
	}

	var from int = (page-1)*pageSize + 1
	var to int = int(math.Min(float64(page*pageSize), float64(count)))

	return fmt.Sprintf("%d - %d", from, to)
}

func (d *datasources) SearchTags(userID, searchKeyword string, limit int) ([]bookmarks.TagsResult, error) {
	return d.bookmarks.SearchTags(userID, searchKeyword, limit)
}

func (d *datasources) GetStatistics(userID string) (Statistics, error) {
	count, err := d.bookmarks.Count(userID)
	stats := Statistics{Count: count}

	return stats, err
}

func (d *datasources) ListBookmarks(userID string, q *string, p, ps *int, sb, sd *string, t *[]string) (BookmarksPage, error) {
	searchKeyword := ""
	if q != nil {
		searchKeyword = *q
	}

	page := 1
	if p != nil {
		page = *p
	}

	pageSize := 10
	if ps != nil {
		pageSize = *ps
	}

	sortBy := "updatedAt"
	if sb != nil {
		sortBy = *sb
	}

	sortDirection := "desc"
	if sd != nil {
		sortDirection = *sd
	}

	tags := []string{}
	if t != nil {
		tags = *t
	}

	bookmarks, count, listErr := d.bookmarks.List(userID, searchKeyword, page, pageSize, sortBy, sortDirection, tags)

	var items []Bookmark = make([]Bookmark, 0, len(bookmarks))
	for _, b := range bookmarks {
		bookmark := Bookmark{
			ShortId:    b.ShortId,
			URL:        b.URL,
			UpdatedAt:  b.UpdatedAt,
			UpdatedBy:  b.UpdatedBy,
			CreatedBy:  b.CreatedBy,
			CreatedAt:  b.CreatedAt,
			Title:      b.Title,
			FaviconURL: b.FaviconURL,
			Tags:       b.Tags,
		}
		items = append(items, bookmark)
	}

	return BookmarksPage{
		Items: items,
		Meta: PaginationMeta{
			Count:         count,
			Query:         searchKeyword,
			Page:          page,
			PageSize:      pageSize,
			SortBy:        sortBy,
			SortDirection: sortDirection,
			Range:         getPaginationRange(count, page, pageSize),
		},
	}, listErr

}

func (d *datasources) CreateBookmark(userID, clientId, URL, title string, tags *[]string) (Bookmark, error) {
	createdBookmark, err := d.bookmarks.InsertOne(userID, URL, title, *tags)

	fmt.Println(createdBookmark)
	go func() {
		_, err1 := d.ParsePageContentOfBookmark(createdBookmark)

		// todo: add logger
		fmt.Println(err1)
	}()

	response := Bookmark{
		ShortId:   createdBookmark.ShortId,
		URL:       createdBookmark.URL,
		UpdatedAt: createdBookmark.UpdatedAt,
		UpdatedBy: createdBookmark.UpdatedBy,
		CreatedBy: createdBookmark.CreatedBy,
		CreatedAt: createdBookmark.CreatedAt,
		Title:     createdBookmark.Title,
		Tags:      createdBookmark.Tags,
	}

	if err == nil {
		go d.mqpublisher.Publish(mqpublisher.BookmarkCreated(userID, clientId, createdBookmark.ShortId, nil))
	}

	return response, err
}

func (d *datasources) UpdateBookmark(bookmarkID, userID, clientId string, update map[string]interface{}) error {
	bookmark, getBookmarkErr := d.bookmarks.GetOne(bookmarkID)
	if getBookmarkErr != nil || bookmark.ShortId == "" {
		return fmt.Errorf("NOT_FOUND")
	}

	if bookmark.CreatedBy != userID {
		return fmt.Errorf("NOT_ALLOWED_TO_UPDATE")
	}

	// only following fields are allowed for update
	allowedFields := []string{
		"title",
		"faviconUrl",
		"tags",
	}
	updatedFields := utils.FilterFieldsOfObject(allowedFields, update)
	err := d.bookmarks.UpdateOne(userID, bookmarkID, updatedFields)

	if err == nil {
		go d.mqpublisher.Publish(mqpublisher.BookmarkUpdated(userID, clientId, bookmarkID, nil))
	}

	return err
}

func (d *datasources) GetBookmark(bookmarkID, userID, clientId string) (Bookmark, error) {
	bookmark, getBookmarkErr := d.bookmarks.GetOne(bookmarkID)

	if userID != bookmark.CreatedBy {
		return Bookmark{}, fmt.Errorf("NOT_ALLOWED_TO_GET")
	}

	creator, errGetCreatorUser := d.users.GetUser(bookmark.CreatedBy)
	if errGetCreatorUser != nil {
		return Bookmark{}, errGetCreatorUser
	}
	updater, errGetUpdaterUser := d.users.GetUser(bookmark.UpdatedBy)
	if errGetUpdaterUser != nil {
		return Bookmark{}, errGetUpdaterUser
	}

	content, getContentErr := d.pageContent.GetLatest(bookmark.URL)
	if getContentErr != nil {
		//
	}

	response := Bookmark{
		ShortId:    bookmark.ShortId,
		URL:        bookmark.URL,
		UpdatedAt:  bookmark.UpdatedAt,
		UpdatedBy:  updater.DisplayName,
		CreatedBy:  creator.DisplayName,
		CreatedAt:  bookmark.CreatedAt,
		Title:      bookmark.Title,
		FaviconURL: bookmark.FaviconURL,
		Tags:       bookmark.Tags,
	}

	go d.mqpublisher.Publish(mqpublisher.BookmarkRead(userID, clientId, bookmarkID, nil))

	return response.AddPageContent(content), getBookmarkErr
}

func (d *datasources) DeleteBookmark(bookmarkID, userID, clientId string) error {
	bookmark, getBookmarkErr := d.bookmarks.GetOne(bookmarkID)
	if getBookmarkErr != nil || bookmark.ShortId == "" {
		return fmt.Errorf("NOT_FOUND")
	}

	if bookmark.CreatedBy != userID {
		return fmt.Errorf("NOT_ALLOWED_TO_DELETE")
	}

	err := d.bookmarks.DeleteOne(bookmarkID)
	if err == nil {
		go d.mqpublisher.Publish(mqpublisher.BookmarkDeleted(userID, clientId, bookmarkID, nil))
	}

	return err
}

func (d *datasources) GetPageContentOfBookmark(bookmarkUrl string) (PageContent, error) {
	pageContent, err := d.pageContent.GetLatest(bookmarkUrl)
	if err != nil || pageContent.Length == 0 {
		return PageContent{}, fmt.Errorf("PAGE_PROPERTIES_NOT_FOUND")
	}

	return PageContent{
		ID:        pageContent.Id,
		URL:       pageContent.URL,
		Title:     pageContent.Title,
		Author:    pageContent.Author,
		Length:    pageContent.Length,
		SiteName:  pageContent.SiteName,
		Image:     pageContent.Image,
		Favicon:   pageContent.Favicon,
		MDContent: pageContent.MDContent,
	}, err
}

func (d *datasources) ParsePageContentOfBookmark(bookmark bookmarks.BookmarkDocument) (PageContent, error) {
	pageContent, scrapErr := d.webScraper.ScrapPage(bookmark.URL)
	if scrapErr != nil {
		return PageContent{}, fmt.Errorf("COULD_NOT_SCRAP_PAGE__%s", &scrapErr)
	}

	ID, _ := gonanoid.Nanoid(8)
	response := PageContent{
		ID:          ID,
		URL:         pageContent.URL,
		Title:       pageContent.Title,
		Author:      pageContent.Author,
		Length:      pageContent.Length,
		Excerpt:     pageContent.Excerpt,
		SiteName:    pageContent.SiteName,
		Image:       pageContent.Image,
		Favicon:     pageContent.Favicon,
		MDContent:   pageContent.MDContent,
		HTMLContent: pageContent.HTMLContent,
	}

	d.pageContent.InsertOne(ID, response.ConvertToDBPageContent())

	updatedFields := make(map[string]interface{})
	updatedFields["faviconUrl"] = pageContent.Favicon

	if bookmark.Title == "" {
		updatedFields["title"] = pageContent.Title
	}

	d.bookmarks.UpdateOne(bookmark.UpdatedBy, bookmark.ShortId, updatedFields)

	return response, nil
}
