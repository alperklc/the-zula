package bookmarksService

import (
	"fmt"
	"testing"
	"time"

	"github.com/alperklc/the-zula/service/infrastructure/db/bookmarks"
	"github.com/alperklc/the-zula/service/infrastructure/logger"
	mqpublisher "github.com/alperklc/the-zula/service/infrastructure/messageQueue/publisher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var l = logger.Get()

func TestBookmark(t *testing.T) {
	t.Run("it returns an array of bookmarks and meta information of the page, when list of bookmarks is requested", func(t *testing.T) {
		// Arrange mocks
		paginationMeta := PaginationMeta{
			Count:         0,
			Query:         "test",
			Page:          1,
			PageSize:      10,
			SortBy:        "title",
			SortDirection: "asc",
			Range:         "",
		}

		bookmarkRepository := new(bookmarks.MockedBookmarks)
		bookmarkRepository.On("List", "userId", "test", 1, 10, "title", "asc", []string{}).Return([]bookmarks.BookmarkDocument(nil), 0, nil)
		bookmarksController := NewService(l, nil, bookmarkRepository, nil, nil, nil)

		q, p, ps, sb, sd, tags := "test", 1, 10, "title", "asc", []string{}

		// Act
		listResponse, err := bookmarksController.ListBookmarks("userId", &q, &p, &ps, &sb, &sd, &tags)

		// Assert
		assert := assert.New(t)
		assert.Nil(err)
		assert.Equal(listResponse.Meta, paginationMeta)
	})
	/*
		t.Run("it publishes a message to queue after the bookmark is successfully created", func(t *testing.T) {
			// Arrange mocks
			rk := messagequeue.RK_NOTI_SCR
			bookmarkMessage := mqpublisher.ActivityMessage{
				RoutingKey:   &rk,
				UserID:       "userId",
				ClientID:     "clientId",
				ResourceType: mqpublisher.ResourceTypeBookmark,
				Action:       mqpublisher.ActionCreate,
				Timestamp:    time.Now(),
			}

			bookmarkRepository := new(bookmarks.MockedBookmarks)
			bookmarkRepository.On("InsertOne", "userId", "google.com", "google", []string{}).Return(bookmarks.BookmarkDocument{ShortId: "test"}, nil)

			eventPublisher := new(mqpublisher.MockedMessagePublisher)
			eventPublisher.On("Publish", bookmarkMessage).Return(nil)

			bookmarksController := NewService(l, nil, bookmarkRepository, nil, nil, eventPublisher)

			// Act
			_, err := bookmarksController.CreateBookmark("userId", "clientId", "google.com", "google", &[]string{})

			// Assert
			assert := assert.New(t)
			assert.Nil(err)
			eventPublisher.AssertCalled(t, "Publish", bookmarkMessage)
		})

		t.Run("it publishes a message to queue after the bookmark is successfully updated", func(t *testing.T) {
			// Arrange mocks
			updates := map[string]interface{}{}
			userId, clientId, shortId := "userId", "clientId", "shortId"

			bookmarkMessage := mqpublisher.BookmarkUpdated(userId, clientId, shortId, nil)

			bookmarkRepository := new(bookmarks.MockedBookmarks)
			bookmarkRepository.On("GetOne", shortId).Return(bookmarks.BookmarkDocument{Id: "bookmarkId", CreatedBy: userId, ShortId: shortId}, nil)
			bookmarkRepository.On("UpdateOne", userId, shortId, updates).Return(nil)

			eventPublisher := new(mqpublisher.MockedMessagePublisher)
			eventPublisher.On("Publish", bookmarkMessage).Return(nil)

			bookmarksController := NewService(l, nil, bookmarkRepository, nil, nil, eventPublisher)

			// Act
			err := bookmarksController.UpdateBookmark(shortId, userId, clientId, updates)

			// Assert
			assert := assert.New(t)
			assert.Nil(err)
			eventPublisher.AssertCalled(t, "Publish", bookmarkMessage)
		})
	*/
	t.Run("it returns NOT_FOUND when bookmark is not found", func(t *testing.T) {
		// Arrange mocks
		bookmarkRepository := new(bookmarks.MockedBookmarks)
		bookmarkRepository.On("GetOne", "bookmarkId").Return(bookmarks.BookmarkDocument{}, fmt.Errorf("NOT_FOUND"))

		bookmarksController := NewService(l, nil, bookmarkRepository, nil, nil, nil)

		// Act
		err := bookmarksController.UpdateBookmark("bookmarkId", "userId", "clientId", nil)

		// Assert
		assert := assert.New(t)
		assert.EqualError(err, "NOT_FOUND")
		bookmarkRepository.AssertCalled(t, "GetOne", "bookmarkId")
	})

	t.Run("it returns NOT_ALLOWED_TO_UPDATE when user did not create the bookmark", func(t *testing.T) {
		// Arrange mocks
		bookmarkRepository := new(bookmarks.MockedBookmarks)
		bookmarkRepository.On("GetOne", "bookmarkId").Return(bookmarks.BookmarkDocument{ShortId: "bookmarkId", CreatedBy: "otherUserId"}, nil)

		bookmarksController := NewService(l, nil, bookmarkRepository, nil, nil, nil)

		// Act
		err := bookmarksController.UpdateBookmark("bookmarkId", "userId", "clientId", nil)

		// Assert
		assert := assert.New(t)
		assert.EqualError(err, "NOT_ALLOWED_TO_UPDATE")
		bookmarkRepository.AssertCalled(t, "GetOne", "bookmarkId")
	})

	t.Run("it deletes a bookmark and publishes a message to the queue", func(t *testing.T) {
		// Arrange mocks
		bookmarkMessage := mqpublisher.BookmarkDeleted("userId", "clientId", "bookmarkId", nil)

		bookmarkRepository := new(bookmarks.MockedBookmarks)
		bookmarkRepository.On("GetOne", "bookmarkId").Return(bookmarks.BookmarkDocument{ShortId: "bookmarkId", CreatedBy: "userId"}, nil)
		bookmarkRepository.On("DeleteOne", "bookmarkId").Return(nil)

		eventPublisher := new(mqpublisher.MockedMessagePublisher)
		eventPublisher.On("Publish", mock.MatchedBy(func(msg mqpublisher.ActivityMessage) bool {
			return msg.UserID == bookmarkMessage.UserID &&
				msg.ClientID == bookmarkMessage.ClientID &&
				msg.Action == bookmarkMessage.Action &&
				msg.ObjectID == bookmarkMessage.ObjectID &&

				// Compare the timestamp with less precision
				msg.Timestamp.Truncate(time.Minute).Equal(bookmarkMessage.Timestamp.Truncate(time.Minute))
		})).Return(nil)

		bookmarksController := NewService(l, nil, bookmarkRepository, nil, nil, eventPublisher)

		// Act
		err := bookmarksController.DeleteBookmark("bookmarkId", "userId", "clientId")

		// Assert
		time.Sleep(100 * time.Millisecond)
		assert := assert.New(t)
		assert.NoError(err)
		eventPublisher.AssertExpectations(t)
	})

	t.Run("it returns NOT_FOUND when bookmark does not exist", func(t *testing.T) {
		// Arrange mocks
		bookmarkRepository := new(bookmarks.MockedBookmarks)
		bookmarkRepository.On("GetOne", "bookmarkId").Return(bookmarks.BookmarkDocument{}, fmt.Errorf("NOT_FOUND"))

		bookmarksController := NewService(l, nil, bookmarkRepository, nil, nil, nil)

		// Act
		err := bookmarksController.DeleteBookmark("bookmarkId", "userId", "clientId")

		// Assert
		assert := assert.New(t)
		assert.EqualError(err, "NOT_FOUND")
		bookmarkRepository.AssertCalled(t, "GetOne", "bookmarkId")
	})

	t.Run("it returns NOT_ALLOWED_TO_DELETE when user did not create the bookmark", func(t *testing.T) {
		// Arrange mocks
		bookmarkRepository := new(bookmarks.MockedBookmarks)
		bookmarkRepository.On("GetOne", "bookmarkId").Return(bookmarks.BookmarkDocument{ShortId: "bookmarkId", CreatedBy: "otherUserId"}, nil)

		bookmarksController := NewService(l, nil, bookmarkRepository, nil, nil, nil)

		// Act
		err := bookmarksController.DeleteBookmark("bookmarkId", "userId", "clientId")

		// Assert
		assert := assert.New(t)
		assert.EqualError(err, "NOT_ALLOWED_TO_DELETE")
		bookmarkRepository.AssertCalled(t, "GetOne", "bookmarkId")
	})
}
