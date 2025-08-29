package mapper

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/entity"
)

func ToPostgreUser(user *domain.User) *entity.UserEntity {
	return &entity.UserEntity{
		ID:        user.ID,
		Name:      user.Name,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Image:     user.Image,
	}
}

func ToDomainUser(user *entity.UserEntity) *domain.User {
	return &domain.User{
		ID:        user.ID,
		Name:      user.Name,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Image:     user.Image,
	}
}
