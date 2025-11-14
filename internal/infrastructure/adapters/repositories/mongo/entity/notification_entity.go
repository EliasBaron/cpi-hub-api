package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notification struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Type       string             `bson:"type"`
	EntityType string             `bson:"entity_type"`
	EntityID   int                `bson:"entity_id"`
	UserID     int                `bson:"user_id"`
	Read       bool               `bson:"read"`
	CreatedAt  primitive.DateTime `bson:"created_at"`
}
