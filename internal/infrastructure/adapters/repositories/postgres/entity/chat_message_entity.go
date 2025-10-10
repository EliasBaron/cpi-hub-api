package entity

import "time"

type ChatMessageEntity struct {
	ID        string    `db:"id"`
	Content   string    `db:"content"`
	UserID    string    `db:"user_id"`
	Username  string    `db:"username"`
	SpaceID   string    `db:"space_id"`
	Timestamp time.Time `db:"timestamp"`
}
