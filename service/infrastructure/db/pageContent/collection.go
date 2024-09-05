package pageContent

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

const collectionName = "page-content"

type Collection interface {
	InsertOne(id string, input PageContent) error
	GetLatest(url string) (PageContentDocument, error)
	ImportMany(pageContent []PageContentDocument) error
}

type db struct {
	collection *mongo.Collection
}

func NewDb(d *mongo.Database) Collection {
	return &db{
		collection: d.Collection(collectionName),
	}
}

func (d *db) InsertOne(id string, input PageContent) error {
	createdAt := time.Now()

	pageContent := PageContentDocument{
		CreatedAt: createdAt,
		URL:       input.URL,
		Title:     input.Title,
		Author:    input.Author,
		Length:    input.Length,
		SiteName:  input.SiteName,
		Image:     input.Image,
		Favicon:   input.Favicon,
		MDContent: input.MDContent,
	}

	_, err := d.collection.InsertOne(context.TODO(), pageContent)

	return err
}

func (d *db) GetLatest(url string) (PageContentDocument, error) {
	filter := bson.M{"url": url}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"_id", -1}})
	findOptions.SetLimit(1)

	cursor, err := d.collection.Find(context.TODO(), &filter, findOptions)
	if err != nil {
		return PageContentDocument{}, err
	}

	var entry PageContentDocument
	var decodeError error
	for cursor.Next(context.TODO()) {
		decodeError = cursor.Decode(&entry)
	}
	if decodeError != nil {
		return PageContentDocument{}, decodeError
	}

	cursor.Close(context.TODO())

	return entry, err
}

func (d *db) ImportMany(pageContent []PageContentDocument) error {
	var itemsToInsert []interface{} = make([]interface{}, 0, len(pageContent))
	for _, pc := range pageContent {
		itemsToInsert = append(itemsToInsert, pc)
	}

	_, err := d.collection.InsertMany(context.TODO(), itemsToInsert)
	return err
}
