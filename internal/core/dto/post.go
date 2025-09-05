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

type PostExtendedDTO struct {
	ID        int                  `json:"id"`
	Title     string               `json:"title"`
	Content   string               `json:"content"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	UpdatedBy int                  `json:"updated_by"`
	CreatedBy UserDTO              `json:"created_by"`
	Space     SimpleSpaceDto       `json:"space"`
	Comments  []CommentWithUserDTO `json:"comments"`
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

func ToPostExtendedDTO(post *domain.ExtendedPost) PostExtendedDTO {
	commentsDTO := make([]CommentWithUserDTO, 0, len(post.Comments))

	for _, c := range post.Comments {
		commentsDTO = append(commentsDTO, ToCommentWithUserAndPostDTO(c))
	}

	return PostExtendedDTO{
		ID:        post.Post.ID,
		Title:     post.Post.Title,
		Content:   post.Post.Content,
		CreatedAt: post.Post.CreatedAt,
		UpdatedAt: post.Post.UpdatedAt,
		UpdatedBy: post.Post.UpdatedBy,
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
		Comments: commentsDTO,
	}
}

func ToPostDTOs(posts []*domain.Post) []PostDTO {
	postDTOs := make([]PostDTO, 0, len(posts))
	for _, post := range posts {
		postDTOs = append(postDTOs, ToPostDTO(post))
	}

	return postDTOs
}

func ToPostExtendedDTOs(posts []*domain.ExtendedPost) []PostExtendedDTO {
	postExtendedDTOs := make([]PostExtendedDTO, 0, len(posts))

	for _, post := range posts {
		postExtendedDTOs = append(postExtendedDTOs, ToPostExtendedDTO(post))
	}

	return postExtendedDTOs
}
