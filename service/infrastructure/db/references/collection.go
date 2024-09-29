package references

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "references"

type Collection interface {
	ListReferencesOfNoteInDepth(noteId string, depth int) ([]ReferencesDocument, error)
	InsertMany(from string, to []string) error
	DeleteAllReferencesFromNote(noteId string) error
	DeleteAllReferencesToNote(noteId string) error
	ImportMany(pageContent []ReferencesDocument) (int, error)
	Export(noteIds []string) ([]ReferencesDocument, error)
}

type db struct {
	collection *mongo.Collection
}

func NewDb(d *mongo.Database) Collection {
	return &db{
		collection: d.Collection(collectionName),
	}
}

func (d *db) ListReferencesOfNoteInDepth(noteId string, depth int) ([]ReferencesDocument, error) {
	pipeline := []bson.M{
		{"$facet": bson.D{{"forwardNodes", bson.A{bson.M{"$match": bson.M{"from": noteId}}}}, {"backwardNodes", bson.A{bson.M{"$match": bson.M{"to": noteId}}}}}},
		{"$graphLookup": bson.M{"from": "notes_references", "startWith": "$forwardNodes.to", "connectToField": "from", "connectFromField": "to", "as": "forwardNodesInDepth", "maxDepth": 1}},
		{"$graphLookup": bson.M{"from": "notes_references", "startWith": "$backwardNodes.from", "connectToField": "to", "connectFromField": "from", "as": "backwardNodesInDepth", "maxDepth": 1}},
		{"$addFields": bson.M{"allItems": bson.M{"$concatArrays": bson.A{"$forwardNodes", "$backwardNodes", "$forwardNodesInDepth", "$backwardNodesInDepth"}}}},
		{"$project": bson.M{"allItems": 1}},
		{"$unwind": bson.M{"path": "$allItems"}},
		{"$replaceRoot": bson.M{"newRoot": "$allItems"}},
	}

	cursor, err := d.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	var aggregationResult []ReferencesDocument
	if decodeError := cursor.All(context.TODO(), &aggregationResult); decodeError != nil {
		return nil, decodeError
	}

	cursor.Close(context.TODO())
	return aggregationResult, nil
}

func (d *db) InsertMany(referenceFrom string, idsOfTargetTotes []string) error {
	var itemsToInsert []interface{} = make([]interface{}, 0, len(idsOfTargetTotes))
	for i := range idsOfTargetTotes {

		noteReference := ReferencesDocument{
			From: referenceFrom,
			To:   idsOfTargetTotes[i],
		}
		itemsToInsert = append(itemsToInsert, noteReference)
	}

	_, err := d.collection.InsertMany(context.TODO(), itemsToInsert)
	return err
}

func (d *db) DeleteAllReferencesFromNote(noteId string) error {
	_, err := d.collection.DeleteMany(context.TODO(), bson.M{"from": noteId})
	return err
}

func (d *db) DeleteAllReferencesToNote(noteId string) error {
	_, err := d.collection.DeleteMany(context.TODO(), bson.M{"to": noteId})
	return err
}

func (d *db) ImportMany(refs []ReferencesDocument) (int, error) {
	var itemsToInsert []interface{} = make([]interface{}, 0, len(refs))
	for _, pc := range refs {
		itemsToInsert = append(itemsToInsert, pc)
	}

	result, err := d.collection.InsertMany(context.TODO(), itemsToInsert)
	return len(result.InsertedIDs), err
}

func (d *db) Export(noteIds []string) ([]ReferencesDocument, error) {
	var referencesDocuments []ReferencesDocument
	filter := bson.M{
		"$or": []bson.M{
			{"from": bson.M{"$in": noteIds}},
			{"to": bson.M{"$in": noteIds}},
		},
	}

	cursor, findErr := d.collection.Find(context.TODO(), filter)
	if findErr != nil {
		return nil, findErr
	}

	if decodeError := cursor.All(context.TODO(), &referencesDocuments); decodeError != nil {
		return nil, decodeError
	}

	defer cursor.Close(context.TODO())

	return referencesDocuments, nil
}
