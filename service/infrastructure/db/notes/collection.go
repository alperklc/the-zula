package notes

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "notes"

type Collection interface {
	SearchTags(userId, searchKeyword string, limit int) ([]TagsResult, error)
	Count(userId string) (int64, error)
	List(userId, searchKeyword string, page, pageSize int, sortBy, sortDirection string, tags []string) ([]NoteDocument, int, error)
	GetNotes(ids, fields []string) ([]NoteDocument, error)
	GetOne(Id string) (NoteDocument, error)
	InsertOne(userId, title, content string, tags []string) (NoteDocument, error)
	UpdateOne(userId, Id string, update interface{}) error
	DeleteOne(Id string) error
}

type db struct {
	collection *mongo.Collection
}

func NewNotesRepository(d *mongo.Database) Collection {
	return &db{
		collection: d.Collection(collectionName),
	}
}

func mapSortDirection(sd string) int {
	if sd == "desc" {
		return -1
	}

	return 1
}

func (d *db) SearchTags(userId, searchKeyword string, limit int) ([]TagsResult, error) {
	matchByUser := bson.D{{
		"$match", bson.M{"createdBy": userId},
	}}

	unwind := bson.D{{
		"$unwind", "$tags",
	}}

	sortByCount := bson.D{{
		"$sortByCount", "$tags",
	}}

	matchByKeyword := bson.D{{
		"$match", bson.M{"_id": bson.M{"$regex": searchKeyword, "$options": "i"}},
	}}

	limitResults := bson.D{{
		"$limit", limit,
	}}

	cursor, err := d.collection.Aggregate(context.Background(), mongo.Pipeline{matchByUser, unwind, sortByCount, matchByKeyword, limitResults})
	if err != nil {
		return nil, err
	}

	aggregationResult := []TagsResult{}
	if decodeError := cursor.All(context.Background(), &aggregationResult); decodeError != nil {
		return nil, decodeError
	}

	defer cursor.Close(context.Background())

	return aggregationResult, nil
}

func (d *db) Count(userId string) (int64, error) {
	filter := bson.M{"createdBy": userId}
	return d.collection.CountDocuments(context.TODO(), filter)
}

func (d *db) List(userId, searchKeyword string, page, pageSize int, sortBy, sortDirection string, tags []string) ([]NoteDocument, int, error) {
	skip := (page - 1) * pageSize

	searchKeywordMatch := []bson.M{
		{"title": bson.M{"$regex": searchKeyword, "$options": "i"}},
		{"tags": bson.M{"$elemMatch": bson.M{"$regex": searchKeyword, "$options": "i"}}},
	}

	matchFilter := bson.M{
		"$or":       searchKeywordMatch,
		"createdBy": userId,
	}

	if len(tags) > 0 {
		matchFilter["tags"] = bson.M{"$in": tags}
	}

	matchStage := bson.D{{
		"$match", matchFilter,
	}}

	facetStage := bson.D{{"$facet",
		bson.D{
			{"count", bson.A{bson.M{"$count": "count"}}},
			{"items", bson.A{
				bson.M{"$sort": bson.M{sortBy: mapSortDirection(sortDirection)}},
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

	cursor, err := d.collection.Aggregate(context.Background(), mongo.Pipeline{matchStage, facetStage, projectStage})
	if err != nil {
		return make([]NoteDocument, 0), 0, err
	}

	var aggregationResult []NotesAggregation

	if decodeError := cursor.All(context.Background(), &aggregationResult); decodeError != nil {
		return make([]NoteDocument, 0), 0, decodeError
	}

	defer cursor.Close(context.Background())

	count := aggregationResult[0].Facets.Count
	return aggregationResult[0].Items, count, nil
}

func (d *db) GetNotes(ids, fields []string) ([]NoteDocument, error) {
	var noteDocuments []NoteDocument
	filter := bson.M{"_id": bson.M{"$in": ids}}

	var projection = make(map[string]interface{})
	for _, field := range fields {
		projection[field] = 1
	}

	findOptions := &options.FindOptions{Projection: projection}
	cursor, findErr := d.collection.Find(context.TODO(), filter, findOptions)
	if findErr != nil {
		return nil, findErr
	}

	if decodeError := cursor.All(context.TODO(), &noteDocuments); decodeError != nil {
		return nil, decodeError
	}

	defer cursor.Close(context.TODO())

	return noteDocuments, nil
}

func (d *db) GetOne(Id string) (NoteDocument, error) {
	oid, _ := primitive.ObjectIDFromHex(Id)

	var noteDocument NoteDocument
	filter := bson.M{"_id": oid}
	err := d.collection.FindOne(context.TODO(), filter).Decode(&noteDocument)

	return noteDocument, err
}

func (d *db) InsertOne(userId, title, content string, tags []string) (NoteDocument, error) {
	createdAt := time.Now()
	noteObject := NoteDocument{
		UpdatedAt: createdAt,
		UpdatedBy: userId,
		CreatedBy: userId,
		CreatedAt: createdAt,
		Title:     title,
		Content:   content,
		Tags:      tags,
	}

	_, err := d.collection.InsertOne(context.TODO(), noteObject)
	return noteObject, err
}

func (d *db) UpdateOne(userId, Id string, updates interface{}) error {
	input, _ := updates.(map[string]interface{})
	input["updatedBy"] = userId
	input["updatedAt"] = time.Now()

	var document bson.M
	obj, marshalErr := bson.Marshal(input)
	if marshalErr != nil {
		return marshalErr
	}

	unmarshalErr := bson.Unmarshal(obj, &document)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	_, err := d.collection.UpdateOne(context.TODO(), bson.M{"_id": Id}, bson.M{"$set": document})
	return err
}

func (d *db) DeleteOne(Id string) error {
	_, err := d.collection.DeleteOne(context.TODO(), bson.M{"_id": Id})

	return err
}
