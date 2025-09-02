package dto

import (
	"cpi-hub-api/internal/core/domain"
	"time"
)

type CreatePost struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	CreatedBy int    `json:"created_by" binding:"required"`
	SpaceID   int    `json:"space_id" binding:"required"`
}

type PostDTO struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy int       `json:"updated_by"`
	SpaceID   int       `json:"space_id"`
}

type PostWithUserSpaceDTO struct {
	ID        int            `json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	CreatedBy UserDTO        `json:"created_by"`
	Space     SimpleSpaceDto `json:"space"`
}

func (c *CreatePost) ToDomain() *domain.Post {
	return &domain.Post{
		Title:     c.Title,
		Content:   c.Content,
		CreatedBy: c.CreatedBy,
		SpaceID:   c.SpaceID,
	}
}

func ToPostDTO(post *domain.Post) PostDTO {
	return PostDTO{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedBy: post.CreatedBy,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		UpdatedBy: post.UpdatedBy,
		SpaceID:   post.SpaceID,
	}
}

func ToPostWithUserSpaceDTO(post *domain.PostWithUserSpace) PostWithUserSpaceDTO {
	return PostWithUserSpaceDTO{
		ID:      post.Post.ID,
		Title:   post.Post.Title,
		Content: post.Post.Content,
		CreatedBy: UserDTO{
			ID:       post.User.ID,
			Name:     post.User.Name,
			LastName: post.User.LastName,
			Image:    post.User.Image,
			Email:    post.User.Email,
		},
		Space: SimpleSpaceDto{
			ID:   post.Space.ID,
			Name: post.Space.Name,
		},
	}
}
