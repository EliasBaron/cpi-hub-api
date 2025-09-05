package dto

import (
	"cpi-hub-api/internal/core/domain"
	"time"
)

type CreateSpace struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	CreatedBy   int    `json:"created_by" binding:"required"`
}

type SpaceDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   int    `json:"created_by"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	UpdatedBy   int    `json:"updated_by"`
}

type SimpleSpaceDto struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SpaceWithUserDTO struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	CreatedBy   UserDTO `json:"created_by"`
}

func (c *CreateSpace) ToDomain() *domain.Space {
	return &domain.Space{
		Name:        c.Name,
		Description: c.Description,
		CreatedBy:   c.CreatedBy,
	}
}

func ToSpaceWithUserDTO(space *domain.SpaceWithUser) SpaceWithUserDTO {
	return SpaceWithUserDTO{
		ID:          space.Space.ID,
		Name:        space.Space.Name,
		Description: space.Space.Description,
		CreatedAt:   space.Space.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   space.Space.UpdatedAt.Format(time.RFC3339),
		CreatedBy: UserDTO{
			ID:       space.User.ID,
			Name:     space.User.Name,
			LastName: space.User.LastName,
			Image:    space.User.Image,
			Email:    space.User.Email,
		},
	}
}

func ToSpaceDTO(space *domain.Space) SpaceDTO {
	return SpaceDTO{
		ID:          space.ID,
		Name:        space.Name,
		Description: space.Description,
		CreatedBy:   space.CreatedBy,
		CreatedAt:   space.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   space.UpdatedAt.Format(time.RFC3339),
		UpdatedBy:   space.UpdatedBy,
	}
}

func ToSpaceWithUserDTOs(spaces []*domain.SpaceWithUser) []SpaceWithUserDTO {
	spacesWithUserDTOs := make([]SpaceWithUserDTO, 0, len(spaces))

	for _, space := range spaces {
		spacesWithUserDTOs = append(spacesWithUserDTOs, ToSpaceWithUserDTO(space))
	}
	return spacesWithUserDTOs
}
