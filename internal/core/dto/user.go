package dto

import (
	"cpi-hub-api/internal/core/domain"
	"time"
)

type CreateUser struct {
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
	Image    string `json:"image"`
}

type UserDTO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Image    string `json:"image"`
}

type UserDTOWithSpaces struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	LastName string     `json:"last_name"`
	Email    string     `json:"email"`
	Image    string     `json:"image"`
	Spaces   []SpaceDTO `json:"spaces"`
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

func ToUserDTOWithSpaces(user *domain.UserWithSpaces) UserDTOWithSpaces {
	spaceDTOs := make([]SpaceDTO, 0, len(user.Spaces))

	for _, s := range user.Spaces {
		spaceDTOs = append(spaceDTOs, SpaceDTO{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			CreatedBy:   user.User.ID,
			CreatedAt:   s.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   s.UpdatedAt.Format(time.RFC3339),
			UpdatedBy:   user.User.ID,
		})
	}

	return UserDTOWithSpaces{
		ID:       user.User.ID,
		Name:     user.User.Name,
		LastName: user.User.LastName,
		Email:    user.User.Email,
		Image:    user.User.Image,
		Spaces:   spaceDTOs,
	}
}

type UpdateUserSpacesDTO struct {
	UserID   int
	SpaceIDs []int
	Action   string
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SearchUsersParams struct {
	FullName      string `form:"full_name"`
	Page          int    `form:"page"`
	PageSize      int    `form:"page_size"`
	OrderBy       string `form:"order_by"`
	SortDirection string `form:"sort_direction"`
}

type PaginatedUsersResponse struct {
	Data     []UserDTO `json:"data"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Total    int       `json:"total"`
}

type UpdateUserDTO struct {
	UserID   int
	Name     *string `json:"name"`
	LastName *string `json:"last_name"`
	Image    *string `json:"image"`
}
