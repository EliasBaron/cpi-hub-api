package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type News struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	Content     string              `bson:"content"`
	Image       string              `bson:"image,omitempty"`
	RedirectURL string              `bson:"redirect_url,omitempty"`
	ExpiresAt   *primitive.DateTime `bson:"expires_at,omitempty"`
	CreatedAt   primitive.DateTime  `bson:"created_at"`
}
