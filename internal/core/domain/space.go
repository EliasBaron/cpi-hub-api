package domain

import "time"

type Space struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
}
