package dto

import "cpi-hub-api/internal/core/domain"

type CreateUser struct {
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
	Image    string `json:"image"`
}

type UserDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Image    string `json:"image"`
}

func (c *CreateUser) ToDomain() *domain.User {
	return &domain.User{
		Name:     c.Name,
		LastName: c.LastName,
		Email:    c.Email,
		Password: c.Password,
		Image:    c.Image,
	}
}

func ToUserDTO(user *domain.User) UserDTO {
	return UserDTO{
		ID:       user.ID,
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
		Image:    user.Image,
	}
}
