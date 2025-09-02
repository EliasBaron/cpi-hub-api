package space

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
	"time"
)

type SpaceUseCase interface {
	Create(ctx context.Context, space *domain.Space) (*domain.SpaceWithUser, error)
	Get(ctx context.Context, id string) (*domain.SpaceWithUser, error)
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

	err = s.userSpaceRepository.AddUserToSpace(ctx, existingUser.ID, space.ID)
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
