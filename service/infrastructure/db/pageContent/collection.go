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
	ImportMany(pageContent []PageContentDocument) (int, error)
	ExportContent(urls []string) ([]PageContentDocument, error)
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

func (d *db) ImportMany(pageContent []PageContentDocument) (int, error) {
	var itemsToInsert []interface{} = make([]interface{}, 0, len(pageContent))
	for _, pc := range pageContent {
		itemsToInsert = append(itemsToInsert, pc)
	}

	result, err := d.collection.InsertMany(context.TODO(), itemsToInsert)
	return len(result.InsertedIDs), err
}

func (d *db) ExportContent(urls []string) ([]PageContentDocument, error) {
	ctx := context.TODO()
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"url", bson.D{{"$in", urls}}}}}},
		{{"$sort", bson.D{{"_id", -1}}}},
		{{"$group", bson.D{
			{"_id", "$url"},
			{"latestDocument", bson.D{{"$first", "$$ROOT"}}},
		}}},
		{{"$replaceRoot", bson.D{{"newRoot", "$latestDocument"}}}},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []PageContentDocument
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
