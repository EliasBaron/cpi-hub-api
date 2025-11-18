package dto

import (
	"cpi-hub-api/internal/core/domain"
	"time"
)

type CreatePost struct {
	Title     string  `json:"title" binding:"required"`
	Content   string  `json:"content" binding:"required"`
	Image     *string `json:"image"`
	CreatedBy int     `json:"created_by" binding:"required"`
	SpaceID   int     `json:"space_id" binding:"required"`
}

type UpdatePost struct {
	PostID  int
	Title   string `json:"title"  `
	Content string `json:"content"`
}

type SearchPostsParams struct {
	Page          int
	PageSize      int
	OrderBy       string
	SortDirection string
	SpaceID       int
	UserID        int
	Query         string
}

type InterestedPostsParams struct {
	Page          int
	PageSize      int
	OrderBy       string
	SortDirection string
	UserID        int // required for interested posts
}

type PaginatedPostsResponse struct {
	Data     []PostExtendedDTO `json:"data"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Total    int               `json:"total"`
}

type PostExtendedDTO struct {
	ID        int                  `json:"id"`
	Title     string               `json:"title"`
	Content   string               `json:"content"`
	Image     *string              `json:"image"`
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
		Image:     c.Image,
		CreatedBy: c.CreatedBy,
		SpaceID:   c.SpaceID,
	}
}

func ToPostExtendedDTO(post *domain.ExtendedPost) PostExtendedDTO {
	// Build nested tree of comments with replies
	commentsDTO := ToCommentWithUserTreeDTOs(post.Comments)

	return PostExtendedDTO{
		ID:        post.Post.ID,
		Title:     post.Post.Title,
		Content:   post.Post.Content,
		Image:     post.Post.Image,
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

func ToPostExtendedDTOs(posts []*domain.ExtendedPost) []PostExtendedDTO {
	postExtendedDTOs := make([]PostExtendedDTO, 0, len(posts))

	for _, post := range posts {
		postExtendedDTOs = append(postExtendedDTOs, ToPostExtendedDTO(post))
	}

	return postExtendedDTOs
}
