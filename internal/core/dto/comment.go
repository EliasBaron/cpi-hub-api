package dto

import (
	"cpi-hub-api/internal/core/domain"
	"time"
)

type CreateComment struct {
	Content   string `json:"content" binding:"required"`
	CreatedBy int    `json:"created_by" binding:"required"`
}

type UpdateCommentDTO struct {
	CommentID int    `json:"comment_id" binding:"required"`
	UserID    int    `json:"user_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

type CommentDTO struct {
	ID              int       `json:"id"`
	PostID          int       `json:"post_id"`
	Content         string    `json:"content"`
	CreatedBy       int       `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	ParentCommentID *int      `json:"parent_comment_id,omitempty"`
}

type CommentWithUserDTO struct {
	ID              int                   `json:"id"`
	PostID          int                   `json:"post_id"`
	Content         string                `json:"content"`
	CreatedBy       UserDTO               `json:"created_by"`
	CreatedAt       time.Time             `json:"created_at"`
	ParentCommentID *int                  `json:"parent_comment_id,omitempty"`
	Replies         []*CommentWithUserDTO `json:"replies,omitempty"`
}

func (c *CommentDTO) ToDomain() *domain.Comment {
	return &domain.Comment{
		PostID:    c.PostID,
		Content:   c.Content,
		CreatedBy: c.CreatedBy,
		ParentID:  c.ParentCommentID,
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
		ID:              comment.Comment.ID,
		PostID:          comment.Comment.PostID,
		Content:         comment.Comment.Content,
		CreatedBy:       ToUserDTO(comment.User),
		CreatedAt:       comment.Comment.CreatedAt,
		ParentCommentID: comment.Comment.ParentID,
		Replies:         []*CommentWithUserDTO{},
	}
}

// ToCommentWithUserTreeDTOs builds a nested tree of comments with replies from a flat list
func ToCommentWithUserTreeDTOs(comments []*domain.CommentWithInfo) []CommentWithUserDTO {
	if len(comments) == 0 {
		return []CommentWithUserDTO{}
	}

	byID := make(map[int]*domain.CommentWithInfo, len(comments))
	children := make(map[int][]int)
	roots := make([]int, 0)

	for _, c := range comments {
		byID[c.Comment.ID] = c
	}

	for _, c := range comments {
		if c.Comment.ParentID == nil {
			roots = append(roots, c.Comment.ID)
			continue
		}
		pid := *c.Comment.ParentID
		children[pid] = append(children[pid], c.Comment.ID)
	}

	// Recursive builder
	var build func(id int) *CommentWithUserDTO
	build = func(id int) *CommentWithUserDTO {
		cwi := byID[id]
		if cwi == nil {
			return nil
		}
		dto := &CommentWithUserDTO{
			ID:              cwi.Comment.ID,
			PostID:          cwi.Comment.PostID,
			Content:         cwi.Comment.Content,
			CreatedBy:       ToUserDTO(cwi.User),
			CreatedAt:       cwi.Comment.CreatedAt,
			ParentCommentID: cwi.Comment.ParentID,
		}
		if kids, ok := children[id]; ok {
			dto.Replies = make([]*CommentWithUserDTO, 0, len(kids))
			for _, childID := range kids {
				if childDTO := build(childID); childDTO != nil {
					dto.Replies = append(dto.Replies, childDTO)
				}
			}
		}
		return dto
	}

	out := make([]CommentWithUserDTO, 0, len(roots))
	for _, rootID := range roots {
		if dto := build(rootID); dto != nil {
			out = append(out, *dto)
		}
	}
	return out
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
