package entity

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reaction struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     int                `bson:"user_id"`
	EntityType string             `bson:"entity_type"`
	EntityID   int                `bson:"entity_id"`
	Action     string             `bson:"action"`
	Timestamp  time.Time          `bson:"timestamp"`
}
