package entity

import "time"

type CommentEntity struct {
	ID        int       `db:"id"`
	PostID    int       `db:"post_id"`
	Content   string    `db:"content"`
	Image     *string   `db:"image"`
	CreatedBy int       `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	ParentID  *int      `db:"parent_comment_id,omitempty"`
}
