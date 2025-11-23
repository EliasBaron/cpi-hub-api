package domain

import "time"

type News struct {
	ID          string
	Content     string
	Image       string
	RedirectURL string
	CreatedAt   time.Time
	ExpiresAt   *time.Time
}
