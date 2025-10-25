package mapper

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/entity"
)

func ToPostgresPost(post *domain.Post) *entity.PostEntity {
	return &entity.PostEntity{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Image:     post.Image,
		CreatedBy: post.CreatedBy,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		UpdatedBy: post.UpdatedBy,
		SpaceID:   post.SpaceID,
	}
}

func ToDomainPost(postEntity *entity.PostEntity) *domain.Post {
	return &domain.Post{
		ID:        postEntity.ID,
		Title:     postEntity.Title,
		Content:   postEntity.Content,
		Image:     postEntity.Image,
		CreatedBy: postEntity.CreatedBy,
		CreatedAt: postEntity.CreatedAt,
		UpdatedAt: postEntity.UpdatedAt,
		UpdatedBy: postEntity.UpdatedBy,
		SpaceID:   postEntity.SpaceID,
	}
}
