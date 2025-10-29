package domain

import (
	"time"
)

type Reaction struct {
	ID         string
	UserID     int
	EntityID   string
	EntityType string
	Liked      bool
	Disliked   bool
	Timestamp  time.Time
}
