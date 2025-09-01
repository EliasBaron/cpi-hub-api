package domain

import "time"

type Post struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int
	UpdatedBy int
	SpaceID   int
	Comments  []Comment
}

type PostWithUserSpace struct {
	Post  *Post
	Space *Space
	User  *User
}

type Comment struct {
	ID        int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int
	UpdatedBy int
}
