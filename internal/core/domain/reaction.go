package domain

import (
	"time"
)

type Reaction struct {
	ID         string
	UserID     int
	EntityType string
	EntityID   int
	Liked      bool
	Disliked   bool
	Timestamp  time.Time
}
