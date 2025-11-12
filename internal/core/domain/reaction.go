package domain

import (
	"time"
)

type Reaction struct {
	ID         string
	UserID     int
	EntityType EntityType
	EntityID   int
	Action     ActionType
	Timestamp  time.Time
}

type EntityType string

const (
	EntityTypePost    EntityType = "post"
	EntityTypeComment EntityType = "comment"
)

func IsValidEntityType(entityType string) bool {
	switch EntityType(entityType) {
	case EntityTypePost, EntityTypeComment:
		return true
	}
	return false
}

type ActionType string

const (
	ActionTypeLike    ActionType = "like"
	ActionTypeDislike ActionType = "dislike"
)

func IsValidActionType(actionType string) bool {
	switch ActionType(actionType) {
	case ActionTypeLike, ActionTypeDislike:
		return true
	}
	return false
}
