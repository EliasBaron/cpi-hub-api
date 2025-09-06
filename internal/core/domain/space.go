package domain

import "time"

type Space struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	CreatedBy   int
	UpdatedAt   time.Time
	UpdatedBy   int
}

type SpaceWithUser struct {
	Space *Space
	User  *User
}
