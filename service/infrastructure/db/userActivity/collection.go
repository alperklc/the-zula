package useractivity

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "users-activity"

type Collection interface {
	List(userID string, page, pageSize int, sortBy, sortDirection string) (UserActivityPage, error)
	InsertOne(userID, resourceType, action, objectID string) (UserActivityDocument, error)
	GetMostVisitedContent(userID string, daysSince int, idsOfExcludedEntries []string) ([]UsageStatisticsEntry, error)
	GetLastVisitedContent(userID string, daysSince int, idsOfExcludedEntries []string) ([]UsageStatisticsEntry, error)
	GetIdsOfDeletedEntries(userID string, daysAgo int) ([]string, error)
	GroupActivitiesByDate(userID string) ([]ActivityGraphEntry, error)
}

type useractivity struct {
	collection *mongo.Collection
}

func NewDb(db *mongo.Database) Collection {
	return &useractivity{
		collection: db.Collection(collectionName),
	}
}

func mapSortDirection(sd string) int {
	if sd == "desc" {
		return -1
	}

	return 1
}

func (db *useractivity) List(userID string, page, pageSize int, sortBy, sortDirection string) (UserActivityPage, error) {
	skip := (page - 1) * pageSize

	matchStage := bson.D{{"$match", bson.M{"userID": userID}}}
	facetStage := bson.D{{"$facet",
		bson.D{
			{"meta", bson.A{bson.M{"$count": "count"}}},
			{"items", bson.A{
				bson.M{"$sort": bson.M{sortBy: mapSortDirection(sortDirection)}},
				bson.M{"$skip": skip},
				bson.M{"$limit": pageSize},
			}},
		},
	}}
	projectStage := bson.D{{"$addFields",
		bson.M{
			"meta": bson.M{"$arrayElemAt": bson.A{"$meta", 0}},
		},
	}}

	cursor, err := db.collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, facetStage, projectStage})
	if err != nil {
		return UserActivityPage{}, err
	}

	var aggregationResult []UserActivityPage
	if decodeError := cursor.All(context.TODO(), &aggregationResult); decodeError != nil {
		return UserActivityPage{}, decodeError
	}

	cursor.Close(context.TODO())

	count := aggregationResult[0].Meta.Count
	return UserActivityPage{
		Items: aggregationResult[0].Items,
		Meta: PaginationMeta{
			Count:         count,
			Page:          page,
			PageSize:      pageSize,
			SortBy:        sortBy,
			SortDirection: sortDirection,
			Range:         getPaginationRange(count, page, pageSize),
		},
	}, err
}

func (db *useractivity) InsertOne(userID, resourceType, action, objectID string) (UserActivityDocument, error) {
	timestamp := time.Now()
	_, err := db.collection.InsertOne(context.TODO(), bson.M{
		"timestamp":    timestamp,
		"userID":       userID,
		"resourceType": resourceType,
		"action":       action,
		"objectID":     objectID,
	})

	return UserActivityDocument{
		UserID:       userID,
		ResourceType: resourceType,
		Action:       action,
		ObjectID:     objectID,
		Timestamp:    timestamp,
	}, err
}

func (db *useractivity) GetMostVisitedContent(userID string, daysSince int, idsOfExcludedEntries []string) ([]UsageStatisticsEntry, error) {
	matchFilter := bson.M{
		"userID":       userID,
		"action":       "READ",
		"resourceType": bson.M{"$in": []string{"NOTE", "FILE", "BOOKMARK"}},
		"timestamp":    bson.M{"$gt": time.Now().AddDate(0, 0, daysSince)},
	}
	if idsOfExcludedEntries != nil {
		matchFilter["objectID"] = bson.M{"$nin": idsOfExcludedEntries}
	}

	matchStage := bson.M{"$match": matchFilter}
	groupStage := bson.M{"$group": bson.M{"_id": "$objectID", "count": bson.M{"$sum": 1}, "resourceType": bson.M{"$first": "$resourceType"}, "timestamp": bson.M{"$last": "$timestamp"}, "objectID": bson.M{"$first": "$objectID"}}}
	sortStage := bson.M{"$sort": bson.M{"count": -1, "timestamp": -1}}
	limitStage := bson.M{"$limit": 5}

	cursor, aggregateErr := db.collection.Aggregate(context.TODO(), []bson.M{matchStage, groupStage, sortStage, limitStage})
	if aggregateErr != nil {
		return nil, aggregateErr
	}

	var result []UsageStatisticsEntry
	if decodeError := cursor.All(context.TODO(), &result); decodeError != nil {
		return nil, decodeError
	}

	cursor.Close(context.TODO())

	return result, nil
}

func (db *useractivity) GetLastVisitedContent(userID string, daysSince int, idsOfExcludedEntries []string) ([]UsageStatisticsEntry, error) {
	matchFilter := bson.M{
		"userID":       userID,
		"action":       "READ",
		"resourceType": bson.M{"$in": []string{"NOTE", "FILE", "BOOKMARK"}},
		"timestamp":    bson.M{"$gt": time.Now().AddDate(0, 0, daysSince)},
	}
	if idsOfExcludedEntries != nil {
		matchFilter["objectID"] = bson.M{"$nin": idsOfExcludedEntries}
	}

	matchStage := bson.M{"$match": matchFilter}
	groupStage := bson.M{"$group": bson.M{"_id": "$objectID", "resourceType": bson.M{"$first": "$resourceType"}, "timestamp": bson.M{"$last": "$timestamp"}, "objectID": bson.M{"$first": "$objectID"}}}
	sortStage := bson.M{"$sort": bson.M{"timestamp": -1}}
	limitStage := bson.M{"$limit": 5}

	cursor, aggregateErr := db.collection.Aggregate(context.TODO(), []bson.M{matchStage, groupStage, sortStage, limitStage})
	if aggregateErr != nil {
		return nil, aggregateErr
	}

	var result []UsageStatisticsEntry
	if decodeError := cursor.All(context.TODO(), &result); decodeError != nil {
		return nil, decodeError
	}
	cursor.Close(context.TODO())

	return result, nil
}

func (db *useractivity) GetIdsOfDeletedEntries(userID string, daysAgo int) ([]string, error) {
	matchFilter := bson.M{
		"userID":    userID,
		"action":    "DELETE",
		"timestamp": bson.M{"$gt": time.Now().AddDate(0, 0, daysAgo)},
	}

	matchStage := bson.M{"$match": matchFilter}

	cursor, aggregateErr := db.collection.Aggregate(context.TODO(), []bson.M{matchStage})
	if aggregateErr != nil {
		return nil, aggregateErr
	}

	var result []map[string]interface{}
	if decodeError := cursor.All(context.TODO(), &result); decodeError != nil {
		return nil, decodeError
	}
	cursor.Close(context.TODO())

	var ids []string = make([]string, len(result))
	for i, item := range result {
		ids[i] = item["objectID"].(string)
	}

	return ids, nil
}

func (db *useractivity) GroupActivitiesByDate(userID string) ([]ActivityGraphEntry, error) {
	matchStage := bson.D{{
		"$match", bson.M{
			"userID":    userID,
			"action":    bson.M{"$in": []string{"UPDATE", "CREATE"}},
			"timestamp": bson.M{"$gt": time.Now().AddDate(0, 0, -365)},
		},
	}}

	addFieldsStage := bson.D{{"$addFields", bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$timestamp"}}}}}

	groupStage := bson.D{{"$group", bson.M{"_id": "$date", "count": bson.M{"$sum": 1}}}}

	sortStage := bson.D{{"$sort", bson.M{"_id": 1}}}

	cursor, err := db.collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, addFieldsStage, groupStage, sortStage})
	if err != nil {
		return nil, err
	}

	var result []ActivityGraphEntry
	if decodeError := cursor.All(context.TODO(), &result); decodeError != nil {
		return nil, decodeError
	}
	cursor.Close(context.TODO())

	return result, nil
}
