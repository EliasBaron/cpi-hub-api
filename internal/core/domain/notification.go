package domain

import "time"

type NotificationType string

const (
	NotificationTypeReaction NotificationType = "reaction"
)

type Notification struct {
	ID         string
	Type       NotificationType
	EntityType EntityType
	EntityID   int
	UserID     int
	Read       bool
	CreatedAt  time.Time
}
