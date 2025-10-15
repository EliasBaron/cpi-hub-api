package entity

import "time"

type ChatMessageEntity struct {
	ID        string    `db:"id"`
	Content   string    `db:"content"`
	UserID    int       `db:"user_id"`
	Username  string    `db:"username"`
	SpaceID   int       `db:"space_id"`
	Timestamp time.Time `db:"timestamp"`
}

type ChatMessageEntityWithUser struct {
	ChatMessageEntity
	Image string `db:"image"`
}
