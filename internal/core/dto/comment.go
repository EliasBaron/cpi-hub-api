package dto

import (
	"cpi-hub-api/internal/core/domain"
	"time"
)

type CreateComment struct {
	Content   string `json:"content" binding:"required"`
	CreatedBy int    `json:"created_by" binding:"required"`
}

type CommentDTO struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	Content   string    `json:"content"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentWithUserDTO struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	Content   string    `json:"content"`
	CreatedBy UserDTO   `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *CommentDTO) ToDomain() *domain.Comment {
	return &domain.Comment{
		PostID:    c.PostID,
		Content:   c.Content,
		CreatedBy: c.CreatedBy,
	}
}

func ToCommentWithUserAndPostDTO(comment *domain.CommentWithUser) CommentWithUserDTO {
	return CommentWithUserDTO{
		ID:        comment.Comment.ID,
		PostID:    comment.Comment.PostID,
		Content:   comment.Comment.Content,
		CreatedBy: ToUserDTO(comment.User),
		CreatedAt: comment.Comment.CreatedAt,
	}
}
