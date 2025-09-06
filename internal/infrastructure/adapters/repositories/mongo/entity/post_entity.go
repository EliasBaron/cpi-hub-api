package entity

import (
	"time"
)

type PostEntity struct {
	ID        string    `bson:"_id"`
	Title     string    `bson:"title"`
	Content   string    `bson:"content"`
	SpaceID   string    `bson:"space_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	CreatedBy int       `bson:"created_by"`
	UpdatedBy int       `bson:"updated_by"`
}
