package notesDrafts

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "notes_drafts"

type Collection interface {
	CheckExistence(ids []string) (map[string]bool, error)
	GetOne(id string) (NoteDraftDocument, error)
	UpsertOne(id, title, content string, tags []string) error
	DeleteOne(id string) error
}

type db struct {
	collection *mongo.Collection
}

func NewNotesDraftsRepository(d *mongo.Database) Collection {
	return &db{
		collection: d.Collection(collectionName),
	}
}

func (d *db) CheckExistence(ids []string) (map[string]bool, error) {
	result := make(map[string]bool)
	var noteDraftDocuments []NoteDraftDocument
	filter := bson.M{"id": bson.M{"$in": ids}}

	cursor, findErr := d.collection.Find(context.TODO(), filter)
	if findErr != nil {
		return result, findErr
	}

	if decodeError := cursor.All(context.TODO(), &noteDraftDocuments); decodeError != nil {
		return result, decodeError
	}

	cursor.Close(context.TODO())

	for _, draft := range noteDraftDocuments {
		result[draft.Id] = true
	}

	return result, nil
}

func (d *db) GetOne(id string) (NoteDraftDocument, error) {
	var noteDocument NoteDraftDocument
	filter := bson.M{"_id": id}
	err := d.collection.FindOne(context.TODO(), filter).Decode(&noteDocument)

	return noteDocument, err
}

func (d *db) UpsertOne(id, title, content string, tags []string) error {
	opts := options.Update().SetUpsert(true)

	filter := bson.M{"id": id}
	noteDraftObject := bson.M{"$set": bson.M{
		"id":      id,
		"title":   title,
		"content": content,
		"tags":    tags,
	}}

	_, err := d.collection.UpdateOne(context.TODO(), filter, noteDraftObject, opts)

	return err
}

func (d *db) DeleteOne(id string) error {
	_, err := d.collection.DeleteOne(context.TODO(), bson.M{"id": id})

	return err
}
