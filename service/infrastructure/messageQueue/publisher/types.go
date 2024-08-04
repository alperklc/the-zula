package mqpublisher

import (
	"encoding/json"
	"fmt"
	"time"

	messagequeue "github.com/alperklc/the-zula/service/infrastructure/messageQueue"
)

const (
	ResourceTypeNote     = "NOTE"
	ResourceTypeBookmark = "BOOKMARK"
	ResourceTypeUser     = "USER"
)

const (
	ActionRead   = "READ"
	ActionCreate = "CREATE"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"
)

type ActivityMessage struct {
	RoutingKey   *string   `json:"routingKey" bson:"routingKey"`
	UserID       string    `json:"userID" bson:"userID"`
	ClientID     string    `json:"clientID" bson:"clientID"`
	ResourceType string    `json:"resourceType" bson:"resourceType"`
	Action       string    `json:"action" bson:"action"`
	Timestamp    time.Time `json:"timestamp" bson:"timestamp"`
	ObjectID     string    `json:"objectID" bson:"objectID"`
	Object       *[]byte   `json:"object" bson:"object"`
}

func newActivityMessage(routingKey *string, userId, clientId, objectId, resourceType, action string, object *[]byte) ActivityMessage {
	return ActivityMessage{
		RoutingKey:   routingKey,
		UserID:       userId,
		ClientID:     clientId,
		ObjectID:     objectId,
		ResourceType: resourceType,
		Action:       action,
		Timestamp:    time.Now(),
		Object:       object,
	}
}

func NoteRead(userId, clientId, objectId string, object *[]byte) ActivityMessage {
	routingKey := messagequeue.RK_ONLY_LOG
	return newActivityMessage(&routingKey, userId, clientId, objectId, ResourceTypeNote, ActionRead, object)
}

func NoteCreated(userId, clientId, objectId string, object *[]byte) ActivityMessage {
	routingKey := messagequeue.RK_NOTI_REF
	return newActivityMessage(&routingKey, userId, clientId, objectId, ResourceTypeNote, ActionCreate, object)
}

func NoteUpdated(userId, clientId, objectId string, object map[string]interface{}) ActivityMessage {
	routingKey := messagequeue.RK_NOTI_REF
	jsonData, err := json.Marshal(object)
	if err != nil {
		fmt.Println("Error converting map to JSON:", err)
	}

	return newActivityMessage(&routingKey, userId, clientId, objectId, ResourceTypeNote, ActionUpdate, &jsonData)
}

func NoteDeleted(userId, clientId, objectId string, object *[]byte) ActivityMessage {
	routingKey := messagequeue.RK_NOTIFICA
	return newActivityMessage(&routingKey, userId, clientId, objectId, ResourceTypeNote, ActionDelete, object)
}

func BookmarkRead(userId, clientId, objectId string, object *[]byte) ActivityMessage {
	routingKey := messagequeue.RK_SCR_ONLY
	return newActivityMessage(&routingKey, userId, clientId, objectId, ResourceTypeBookmark, ActionRead, object)
}

func BookmarkCreated(userId, clientId, objectId string, object *[]byte) ActivityMessage {
	routingKey := messagequeue.RK_NOTI_SCR
	return newActivityMessage(&routingKey, userId, clientId, objectId, ResourceTypeBookmark, ActionCreate, object)
}

func BookmarkUpdated(userId, clientId, objectId string, object *[]byte) ActivityMessage {
	routingKey := messagequeue.RK_NOTIFICA
	return newActivityMessage(&routingKey, userId, clientId, objectId, ResourceTypeBookmark, ActionUpdate, object)
}

func BookmarkDeleted(userId, clientId, objectId string, object *[]byte) ActivityMessage {
	routingKey := messagequeue.RK_NOTIFICA
	return newActivityMessage(&routingKey, userId, clientId, objectId, ResourceTypeBookmark, ActionDelete, object)
}

func UserUpdated(userId, clientId string, object *[]byte) ActivityMessage {
	routingKey := messagequeue.RK_NOTI_USR
	return newActivityMessage(&routingKey, userId, clientId, userId, ResourceTypeUser, ActionUpdate, object)
}
