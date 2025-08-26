package dto

import "cpi-hub-api/internal/core/domain"

type CreateSpace struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	CreatedBy   string `json:"created_by" binding:"required"`
}

type SpaceDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   string `json:"created_by"`
	UpdatedAt   string `json:"updated_at"`
	UpdatedBy   string `json:"updated_by"`
}

func (c *CreateSpace) ToDomain() *domain.Space {
	return &domain.Space{
		Name:        c.Name,
		Description: c.Description,
		CreatedBy:   c.CreatedBy,
	}
}

func ToSpaceDTO(space *domain.Space) SpaceDTO {
	return SpaceDTO{
		ID:          space.ID,
		Name:        space.Name,
		Description: space.Description,
		CreatedBy:   space.CreatedBy,
	}
}
