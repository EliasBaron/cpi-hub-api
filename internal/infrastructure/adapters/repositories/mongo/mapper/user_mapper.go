package mapper

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/entity"
)

func ToMongoDatabaseUser(user *domain.User) *entity.UserEntity {
	return &entity.UserEntity{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Image:     user.Image,
		Name:      user.Name,
		LastName:  user.LastName,
	}
}

func ToDomainUser(user *entity.UserEntity) *domain.User {
	return &domain.User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Image:     user.Image,
		Name:      user.Name,
		LastName:  user.LastName,
	}
}
