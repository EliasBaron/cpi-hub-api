package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notification struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	URL         *string            `bson:"url,omitempty"`
	To          int                `bson:"to"`
	Read        bool               `bson:"read"`
	CreatedAt   primitive.DateTime `bson:"created_at"`
}
