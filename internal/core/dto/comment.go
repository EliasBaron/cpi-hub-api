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

type SearchCommentsParams struct {
	Page          int
	PageSize      int
	OrderBy       string
	SortDirection string
	UserID        *int
	PostID        *int
}

type PaginatedCommentsResponse struct {
	Data     []CommentWithSpaceDTO `json:"data"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
	Total    int                   `json:"total"`
}

type CommentWithSpaceDTO struct {
	ID        int            `json:"id"`
	PostID    int            `json:"post_id"`
	Content   string         `json:"content"`
	CreatedBy UserDTO        `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	Space     SimpleSpaceDto `json:"space"`
}

func ToCommentWithUserAndPostDTO(comment *domain.CommentWithInfo) CommentWithUserDTO {
	return CommentWithUserDTO{
		ID:        comment.Comment.ID,
		PostID:    comment.Comment.PostID,
		Content:   comment.Comment.Content,
		CreatedBy: ToUserDTO(comment.User),
		CreatedAt: comment.Comment.CreatedAt,
	}
}

func ToCommentWithSpaceDTO(comment *domain.CommentWithInfo) CommentWithSpaceDTO {
	return CommentWithSpaceDTO{
		ID:        comment.Comment.ID,
		PostID:    comment.Comment.PostID,
		Content:   comment.Comment.Content,
		CreatedBy: ToUserDTO(comment.User),
		CreatedAt: comment.Comment.CreatedAt,
		Space: SimpleSpaceDto{
			ID:   comment.Space.ID,
			Name: comment.Space.Name,
		},
	}
}

func ToCommentWithSpaceDTOs(comments []*domain.CommentWithInfo) []CommentWithSpaceDTO {
	commentDTOs := make([]CommentWithSpaceDTO, 0, len(comments))

	for _, comment := range comments {
		commentDTOs = append(commentDTOs, ToCommentWithSpaceDTO(comment))
	}

	return commentDTOs
}
