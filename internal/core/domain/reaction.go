package domain

import (
	"time"
)

type Reaction struct {
	ID         string
	UserID     int
	EntityType EntityType
	EntityID   int
	Liked      bool
	Disliked   bool
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
