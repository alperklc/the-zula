package notesChanges

import (
	"context"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "notes_changes"

type Collection interface {
	ListHistoryOfNote(noteId string, page, pageSize int) ([]NotesChangesDocument, int, error)
	GetCountOfChanges(noteId string) (int64, error)
	GetOne(noteId, timestamp string) (NotesChangesDocument, error)
	InsertOne(noteId string, updatedAt time.Time, updatedBy, change string) error
}

type db struct {
	collection *mongo.Collection
}

func NewDb(d *mongo.Database) Collection {
	return &db{
		collection: d.Collection(collectionName),
	}
}

func (d *db) ListHistoryOfNote(noteId string, page, pageSize int) ([]NotesChangesDocument, int, error) {
	skip := (page - 1) * pageSize

	matchFilter := bson.M{
		"noteId": noteId,
	}

	matchStage := bson.D{{
		"$match", matchFilter,
	}}

	facetStage := bson.D{{"$facet",
		bson.D{
			{"count", bson.A{bson.M{"$count": "count"}}},
			{"items", bson.A{
				bson.M{"$sort": bson.M{"updatedAt": -1}},
				bson.M{"$skip": skip},
				bson.M{"$limit": pageSize},
			}},
		},
	}}
	projectStage := bson.D{{"$addFields",
		bson.M{
			"facets": bson.M{"$arrayElemAt": bson.A{"$count", 0}},
		},
	}}

	cursor, err := d.collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, facetStage, projectStage})
	if err != nil {
		return make([]NotesChangesDocument, 0), 0, err
	}

	var aggregationResult []NotesChangesAggregation

	if decodeError := cursor.All(context.TODO(), &aggregationResult); decodeError != nil {
		return make([]NotesChangesDocument, 0), 0, decodeError
	}

	cursor.Close(context.TODO())

	count := aggregationResult[0].Facets.Count
	return aggregationResult[0].Items, count, nil
}

func (d *db) GetCountOfChanges(noteId string) (int64, error) {
	filter := bson.M{"noteId": noteId}
	return d.collection.CountDocuments(context.TODO(), filter)
}

func (d *db) GetOne(noteId, timestamp string) (NotesChangesDocument, error) {
	var noteHistoryDocument NotesChangesDocument
	filter := bson.M{"noteId": noteId, "updatedAt": timestamp}
	err := d.collection.FindOne(context.TODO(), filter).Decode(&noteHistoryDocument)

	if err != nil {
		return NotesChangesDocument{}, err
	}

	return noteHistoryDocument, nil

}

func (d *db) InsertOne(noteId string, updatedAt time.Time, updatedBy, change string) error {
	shortId, _ := gonanoid.Nanoid(8)

	noteObject := NotesChangesDocument{
		NoteId:    noteId,
		ShortId:   shortId,
		UpdatedAt: updatedAt,
		UpdatedBy: updatedBy,
		Change:    change,
	}

	_, err := d.collection.InsertOne(context.TODO(), noteObject)

	return err
}
