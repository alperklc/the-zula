package useractivity

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "users-activity"

type Collection interface {
	List(userUID string, page, pageSize int, sortBy, sortDirection string) (UserActivityPage, error)
	InsertOne(userUID, resourceType, action, objectUID string) (UserActivityDocument, error)
	GetMostVisitedContent(userUID string, daysSince int, uidsOfExcludedEntries []string) ([]UsageStatisticsEntry, error)
	GetLastVisitedContent(userUID string, daysSince int, uidsOfExcludedEntries []string) ([]UsageStatisticsEntry, error)
	GetUidsOfDeletedEntries(userUID string, daysAgo int) ([]string, error)
	GroupActivitiesByDate(userUID string) ([]ActivityGraphEntry, error)
}

type useractivity struct {
	logger     *log.Entry
	collection *mongo.Collection
}

func NewUserActivityRepository(db *mongo.Database) Collection {
	return &useractivity{
		logger:     log.WithFields(log.Fields{"package": "database_user_activity"}),
		collection: db.Collection(collectionName),
	}
}

func mapSortDirection(sd string) int {
	if sd == "desc" {
		return -1
	}

	return 1
}

func (db *useractivity) List(userUID string, page, pageSize int, sortBy, sortDirection string) (UserActivityPage, error) {
	db.logger.WithFields(log.Fields{"user": userUID}).Debug("Listing activity records of user")

	skip := (page - 1) * pageSize

	matchStage := bson.D{{"$match", bson.M{"userUID": userUID}}}
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
		db.logger.WithFields(log.Fields{"user": userUID}).Error("Aggregate failed ", err.Error())
		return UserActivityPage{}, err
	}

	var aggregationResult []UserActivityPage
	if decodeError := cursor.All(context.TODO(), &aggregationResult); decodeError != nil {
		db.logger.WithFields(log.Fields{"user": userUID}).Error("Decode failed ", decodeError.Error())
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

func (db *useractivity) InsertOne(userUID, resourceType, action, objectUID string) (UserActivityDocument, error) {
	db.logger.WithFields(log.Fields{"user": userUID}).Debug("Inserting an activity record")

	timestamp := time.Now()
	_, err := db.collection.InsertOne(context.TODO(), bson.M{
		"timestamp":    timestamp,
		"userUID":      userUID,
		"resourceType": resourceType,
		"action":       action,
		"objectUID":    objectUID,
	})

	return UserActivityDocument{
		UserUID:      userUID,
		ResourceType: resourceType,
		Action:       action,
		ObjectUID:    objectUID,
		Timestamp:    timestamp,
	}, err
}

func (db *useractivity) GetMostVisitedContent(userUID string, daysSince int, uidsOfExcludedEntries []string) ([]UsageStatisticsEntry, error) {
	db.logger.WithFields(log.Fields{"user": userUID}).Debug("Getting most visited content")

	matchFilter := bson.M{
		"userUID":      userUID,
		"action":       "READ",
		"resourceType": bson.M{"$in": []string{"NOTE", "FILE", "BOOKMARK"}},
		"timestamp":    bson.M{"$gt": time.Now().AddDate(0, 0, daysSince)},
	}
	if uidsOfExcludedEntries != nil {
		matchFilter["objectUID"] = bson.M{"$nin": uidsOfExcludedEntries}
	}

	matchStage := bson.M{"$match": matchFilter}
	groupStage := bson.M{"$group": bson.M{"_id": "$objectUID", "count": bson.M{"$sum": 1}, "resourceType": bson.M{"$first": "$resourceType"}, "timestamp": bson.M{"$last": "$timestamp"}, "objectUID": bson.M{"$first": "$objectUID"}}}
	sortStage := bson.M{"$sort": bson.M{"count": -1, "timestamp": -1}}
	limitStage := bson.M{"$limit": 5}

	cursor, aggregateErr := db.collection.Aggregate(context.TODO(), []bson.M{matchStage, groupStage, sortStage, limitStage})
	if aggregateErr != nil {
		db.logger.WithFields(log.Fields{"user": userUID}).Error("Aggregate failed ", aggregateErr.Error())
		return nil, aggregateErr
	}

	var result []UsageStatisticsEntry
	if decodeError := cursor.All(context.TODO(), &result); decodeError != nil {
		db.logger.WithFields(log.Fields{"user": userUID}).Error("Decode failed ", decodeError.Error())
		return nil, decodeError
	}

	cursor.Close(context.TODO())

	return result, nil
}

func (db *useractivity) GetLastVisitedContent(userUID string, daysSince int, uidsOfExcludedEntries []string) ([]UsageStatisticsEntry, error) {
	db.logger.WithFields(log.Fields{"user": userUID}).Debug("Getting last visited content")

	matchFilter := bson.M{
		"userUID":      userUID,
		"action":       "READ",
		"resourceType": bson.M{"$in": []string{"NOTE", "FILE", "BOOKMARK"}},
		"timestamp":    bson.M{"$gt": time.Now().AddDate(0, 0, daysSince)},
	}
	if uidsOfExcludedEntries != nil {
		matchFilter["objectUID"] = bson.M{"$nin": uidsOfExcludedEntries}
	}

	matchStage := bson.M{"$match": matchFilter}
	groupStage := bson.M{"$group": bson.M{"_id": "$objectUID", "resourceType": bson.M{"$first": "$resourceType"}, "timestamp": bson.M{"$last": "$timestamp"}, "objectUID": bson.M{"$first": "$objectUID"}}}
	sortStage := bson.M{"$sort": bson.M{"timestamp": -1}}
	limitStage := bson.M{"$limit": 5}

	cursor, aggregateErr := db.collection.Aggregate(context.TODO(), []bson.M{matchStage, groupStage, sortStage, limitStage})
	if aggregateErr != nil {
		db.logger.WithFields(log.Fields{"user": userUID}).Error("Aggregate failed ", aggregateErr.Error())
		return nil, aggregateErr
	}

	var result []UsageStatisticsEntry
	if decodeError := cursor.All(context.TODO(), &result); decodeError != nil {
		db.logger.WithFields(log.Fields{"user": userUID}).Error("Decode failed ", decodeError.Error())
		return nil, decodeError
	}
	cursor.Close(context.TODO())

	return result, nil
}

func (db *useractivity) GetUidsOfDeletedEntries(userUID string, daysAgo int) ([]string, error) {
	db.logger.WithFields(log.Fields{"user": userUID}).Debug("Getting UIDs of deleted entries")

	matchFilter := bson.M{
		"userUID":   userUID,
		"action":    "DELETE",
		"timestamp": bson.M{"$gt": time.Now().AddDate(0, 0, daysAgo)},
	}

	matchStage := bson.M{"$match": matchFilter}

	cursor, aggregateErr := db.collection.Aggregate(context.TODO(), []bson.M{matchStage})
	if aggregateErr != nil {
		db.logger.WithFields(log.Fields{"user": userUID}).Error("Aggregate failed ", aggregateErr.Error())
		return nil, aggregateErr
	}

	var result []map[string]interface{}
	if decodeError := cursor.All(context.TODO(), &result); decodeError != nil {
		db.logger.WithFields(log.Fields{"user": userUID}).Error("Decode failed ", decodeError.Error())
		return nil, decodeError
	}
	cursor.Close(context.TODO())

	var uids []string = make([]string, len(result))
	for i, item := range result {
		uids[i] = item["objectUID"].(string)
	}

	return uids, nil
}

func (db *useractivity) GroupActivitiesByDate(userUID string) ([]ActivityGraphEntry, error) {
	db.logger.WithFields(log.Fields{"user": userUID}).Debug("Grouping entries by date")

	matchStage := bson.D{{
		"$match", bson.M{
			"userUID":   userUID,
			"action":    bson.M{"$in": []string{"UPDATE", "CREATE"}},
			"timestamp": bson.M{"$gt": time.Now().AddDate(0, 0, -365)},
		},
	}}

	addFieldsStage := bson.D{{"$addFields", bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$timestamp"}}}}}

	groupStage := bson.D{{"$group", bson.M{"_id": "$date", "count": bson.M{"$sum": 1}}}}

	sortStage := bson.D{{"$sort", bson.M{"_id": 1}}}

	cursor, err := db.collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, addFieldsStage, groupStage, sortStage})
	if err != nil {
		db.logger.WithFields(log.Fields{"user": userUID}).Error("Aggregate failed ", err.Error())
		return nil, err
	}

	var result []ActivityGraphEntry
	if decodeError := cursor.All(context.TODO(), &result); decodeError != nil {
		db.logger.WithFields(log.Fields{"user": userUID}).Error("Decode failed ", decodeError.Error())
		return nil, decodeError
	}
	cursor.Close(context.TODO())

	return result, nil
}
