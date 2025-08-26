package space

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
	"time"
)

type SpaceUseCase interface {
	Create(ctx context.Context, space *domain.Space) (*domain.Space, error)
	Get(ctx context.Context, id string) (*domain.Space, error)
}

type spaceUseCase struct {
	spaceRepository domain.SpaceRepository
	userRepository  domain.UserRepository
}

func NewSpaceUsecase(spaceRepository domain.SpaceRepository, userRepository domain.UserRepository) SpaceUseCase {
	return &spaceUseCase{
		spaceRepository: spaceRepository,
		userRepository:  userRepository,
	}
}

// Create implements SpaceUseCase.
func (s *spaceUseCase) Create(ctx context.Context, space *domain.Space) (*domain.Space, error) {
	existingUser, err := s.userRepository.FindById(ctx, space.CreatedBy)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, apperror.NewNotFound("User not found", nil, "space_usecase.go:Create")
	}

	existingSpace, err := s.spaceRepository.FindByName(ctx, space.Name)
	if err != nil {
		return nil, err
	}
	if existingSpace != nil {
		return nil, apperror.NewInvalidData("Space with this name already exists", nil, "space_usecase.go:Create")
	}

	space.ID = helpers.NewULID()
	space.CreatedAt = time.Now()
	space.UpdatedBy = space.CreatedBy
	space.UpdatedAt = time.Now()

	err = s.spaceRepository.Create(ctx, space)
	if err != nil {
		return nil, err
	}

	return space, nil
}

// Get implements SpaceUseCase.
func (s *spaceUseCase) Get(ctx context.Context, id string) (*domain.Space, error) {
	space, err := s.spaceRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if space == nil {
		return nil, apperror.NewNotFound("Space not found", nil, "space_usecase.go:Get")
	}

	return space, nil
}
