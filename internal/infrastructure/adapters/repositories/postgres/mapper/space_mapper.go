package mapper

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/entity"
)

func ToPostgresSpace(space *domain.Space) *entity.SpaceEntity {
	return &entity.SpaceEntity{
		ID:          space.ID,
		Name:        space.Name,
		Description: space.Description,
		Members:     space.Members,
		Posts:       space.Posts,
		CreatedBy:   space.CreatedBy,
		CreatedAt:   space.CreatedAt,
		UpdatedBy:   space.UpdatedBy,
		UpdatedAt:   space.UpdatedAt,
	}
}

func ToDomainSpace(space *entity.SpaceEntity) *domain.Space {
	return &domain.Space{
		ID:          space.ID,
		Name:        space.Name,
		Description: space.Description,
		Members:     space.Members,
		Posts:       space.Posts,
		CreatedAt:   space.CreatedAt,
		CreatedBy:   space.CreatedBy,
		UpdatedAt:   space.UpdatedAt,
		UpdatedBy:   space.UpdatedBy,
	}
}
