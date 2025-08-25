package domain

import "time"

type Post struct {
	ID        string
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
	SpaceID   string
	Comments  []Comment
}

type Comment struct {
	ID        string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
}
