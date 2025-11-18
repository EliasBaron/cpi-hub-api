package domain

import "time"

type Post struct {
	ID        int
	Title     string
	Content   string
	Image     *string
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
	Comments []*CommentWithInfo
}

type Comment struct {
	ID        int
	PostID    int
	Content   string
	Image     *string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int
	ParentID  *int
	Replies   []*Comment
}

type CommentWithInfo struct {
	Comment *Comment
	User    *User
	Space   *Space
}
