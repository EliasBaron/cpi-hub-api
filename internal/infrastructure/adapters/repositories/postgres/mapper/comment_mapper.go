package mapper

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/entity"
)

func ToPostgreComment(comment *domain.Comment) *entity.CommentEntity {
	return &entity.CommentEntity{
		ID:        comment.ID,
		PostID:    comment.PostID,
		Content:   comment.Content,
		Image:     comment.Image,
		CreatedBy: comment.CreatedBy,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		ParentID:  comment.ParentID,
	}
}

func ToDomainComment(commentEntity *entity.CommentEntity) *domain.Comment {
	return &domain.Comment{
		ID:        commentEntity.ID,
		PostID:    commentEntity.PostID,
		Content:   commentEntity.Content,
		Image:     commentEntity.Image,
		CreatedBy: commentEntity.CreatedBy,
		CreatedAt: commentEntity.CreatedAt,
		UpdatedAt: commentEntity.UpdatedAt,
		ParentID:  commentEntity.ParentID,
	}
}
