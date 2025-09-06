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

type ExtendedPost struct {
	Post     *Post
	Space    *Space
	User     *User
	Comments []*CommentWithUser
}

type Comment struct {
	ID        int
	PostID    int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int
	UpdatedBy int
}

type CommentWithUser struct {
	Comment *Comment
	User    *User
}
