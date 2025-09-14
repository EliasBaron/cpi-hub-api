package space

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/helpers"
	"cpi-hub-api/pkg/apperror"
	"time"
)

type SpaceUseCase interface {
	Create(ctx context.Context, space *domain.Space) (*domain.SpaceWithUser, error)
	Get(ctx context.Context, id string) (*domain.SpaceWithUser, error)
	Search(ctx context.Context, criteria *domain.SpaceSearchCriteria) (*domain.SearchResult, error)
}

type spaceUseCase struct {
	spaceRepository     domain.SpaceRepository
	userRepository      domain.UserRepository
	userSpaceRepository domain.UserSpaceRepository
}

func NewSpaceUsecase(spaceRepository domain.SpaceRepository, userRepository domain.UserRepository, userSpaceRepository domain.UserSpaceRepository) SpaceUseCase {
	return &spaceUseCase{
		spaceRepository:     spaceRepository,
		userRepository:      userRepository,
		userSpaceRepository: userSpaceRepository,
	}
}

func (s *spaceUseCase) makeSpacesWithUsers(ctx context.Context, spaces []*domain.Space) ([]*domain.SpaceWithUser, error) {
	var spacesWithUsers []*domain.SpaceWithUser

	for _, space := range spaces {
		user, err := helpers.FindEntity(ctx, s.userRepository, "id", space.CreatedBy, "User not found")
		if err != nil {
			return nil, err
		}

		spacesWithUsers = append(spacesWithUsers, &domain.SpaceWithUser{
			Space: space,
			User:  user,
		})
	}

	return spacesWithUsers, nil
}

func (s *spaceUseCase) Create(ctx context.Context, space *domain.Space) (*domain.SpaceWithUser, error) {
	existingUser, err := s.userRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    space.CreatedBy,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, apperror.NewNotFound("User not found", nil, "space_usecase.go:Create")
	}

	existingSpace, err := s.spaceRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "name",
				Value:    space.Name,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if existingSpace != nil {
		return nil, apperror.NewInvalidData("Space with this name already exists", nil, "space_usecase.go:Create")
	}

	space.CreatedAt, space.UpdatedAt = time.Now(), time.Now()
	space.UpdatedBy, space.CreatedBy = existingUser.ID, existingUser.ID

	err = s.spaceRepository.Create(ctx, space)
	if err != nil {
		return nil, err
	}

	err = s.userSpaceRepository.Update(ctx, existingUser.ID, []int{space.ID}, domain.AddUserToSpace)
	if err != nil {
		return nil, err
	}

	return &domain.SpaceWithUser{
		Space: space,
		User:  existingUser,
	}, nil
}

func (s *spaceUseCase) Get(ctx context.Context, id string) (*domain.SpaceWithUser, error) {
	space, err := helpers.FindEntity(ctx, s.spaceRepository, "id", id, "Space not found")
	if err != nil {
		return nil, err
	}

	user, err := helpers.FindEntity(ctx, s.userRepository, "id", space.CreatedBy, "User not found")
	if err != nil {
		return nil, err
	}

	return &domain.SpaceWithUser{
		Space: &domain.Space{
			ID:          space.ID,
			Name:        space.Name,
			Description: space.Description,
			CreatedAt:   space.CreatedAt,
			UpdatedAt:   space.UpdatedAt,
			CreatedBy:   space.CreatedBy,
			UpdatedBy:   space.UpdatedBy,
		},
		User: user,
	}, nil
}

func (s *spaceUseCase) Search(ctx context.Context, searchCriteria *domain.SpaceSearchCriteria) (*domain.SearchResult, error) {
	var direction criteria.Direction
	if searchCriteria.SortDirection == "asc" {
		direction = criteria.OrderDirectionAsc
	} else {
		direction = criteria.OrderDirectionDesc
	}

	criteriaBuilder := criteria.NewCriteriaBuilder().
		WithSort(searchCriteria.OrderBy, direction).
		WithPagination(searchCriteria.Page, searchCriteria.PageSize)

	if searchCriteria.CreatedBy != nil {
		criteriaBuilder.WithFilter("created_by", *searchCriteria.CreatedBy, criteria.OperatorEqual)
	}

	criteria := criteriaBuilder.Build()

	totalCount, err := s.spaceRepository.Count(ctx, criteria)
	if err != nil {
		return nil, err
	}

	spaces, err := s.spaceRepository.FindAll(ctx, criteria)
	if err != nil {
		return nil, err
	}

	var spacesWithUsers []*domain.SpaceWithUser
	if len(spaces) > 0 {
		spacesWithUsers, err = s.makeSpacesWithUsers(ctx, spaces)
		if err != nil {
			return nil, err
		}
	}

	return &domain.SearchResult{
		Data:     spacesWithUsers,
		Page:     searchCriteria.Page,
		PageSize: searchCriteria.PageSize,
		Total:    totalCount,
	}, nil
}
