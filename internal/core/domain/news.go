package domain

import "time"

type News struct {
	ID        int
	Content   string
	Image     string
	CreatedAt time.Time
	ExpiresAt *time.Time
}
