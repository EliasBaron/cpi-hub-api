package entity

import "time"

type SpaceEntity struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	CreatedAt   time.Time `bson:"created_at"`
	CreatedBy   int       `bson:"created_by"`
	UpdatedAt   time.Time `bson:"updated_at"`
	UpdatedBy   int       `bson:"updated_by"`
}
