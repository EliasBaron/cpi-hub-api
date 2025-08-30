package mapper

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/entity"
)

func ToMongoDatabasePost(post *domain.Post) *entity.PostEntity {
	return &entity.PostEntity{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		SpaceID:   post.SpaceID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		CreatedBy: post.CreatedBy,
		UpdatedBy: post.UpdatedBy,
	}
}

func ToDomainPost(postEntity *entity.PostEntity) *domain.Post {
	return &domain.Post{
		ID:        postEntity.ID,
		Title:     postEntity.Title,
		Content:   postEntity.Content,
		SpaceID:   postEntity.SpaceID,
		CreatedAt: postEntity.CreatedAt,
		UpdatedAt: postEntity.UpdatedAt,
		CreatedBy: postEntity.CreatedBy,
		UpdatedBy: postEntity.UpdatedBy,
	}
}
