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

type SearchSpacesDTO struct {
	Name          *string `form:"name" json:"name"`
	CreatedBy     *int    `form:"created_by" json:"created_by"`
	OrderBy       string  `form:"order_by" json:"order_by"`
	Page          int     `form:"page" json:"page"`
	PageSize      int     `form:"page_size" json:"page_size"`
	SortDirection string  `form:"sort_direction" json:"sort_direction"`
}

type SearchSpacesResponseDTO struct {
	Data     []SpaceWithUserAndCountDTO `json:"data"`
	Page     int                        `json:"page"`
	PageSize int                        `json:"page_size"`
	Total    int                        `json:"total"`
}

type SpaceWithUserAndCountDTO struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Users       int     `json:"users"`
	Posts       int     `json:"posts"`
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

func ToSpaceWithUserDTO(space *domain.SpaceWithUserAndCounts) SpaceWithUserAndCountDTO {
	return SpaceWithUserAndCountDTO{
		ID:          space.Space.ID,
		Name:        space.Space.Name,
		Description: space.Space.Description,
		Users:       space.SpaceCounts.Users,
		Posts:       space.SpaceCounts.Posts,
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

func ToSpaceWithUserDTOs(spaces []*domain.SpaceWithUserAndCounts) []SpaceWithUserAndCountDTO {
	spacesWithUserDTOs := make([]SpaceWithUserAndCountDTO, 0, len(spaces))

	for _, space := range spaces {
		spacesWithUserDTOs = append(spacesWithUserDTOs, ToSpaceWithUserDTO(space))
	}
	return spacesWithUserDTOs
}
