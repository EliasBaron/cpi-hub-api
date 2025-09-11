package user

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/pkg/apperror"
	"time"
)

type UserUseCase interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Get(ctx context.Context, id int) (*domain.UserWithSpaces, error)
	Update(ctx context.Context, dto dto.UpdateUserSpacesDTO) error
	GetSpacesByUser(ctx context.Context, userId int) ([]*domain.Space, error)
}

type useCase struct {
	userRepository      domain.UserRepository
	spaceRepository     domain.SpaceRepository
	userSpaceRepository domain.UserSpaceRepository
}

func NewUserUsecase(userRepository domain.UserRepository, spaceRepository domain.SpaceRepository, userSpaceRepository domain.UserSpaceRepository) UserUseCase {
	return &useCase{
		userRepository:      userRepository,
		spaceRepository:     spaceRepository,
		userSpaceRepository: userSpaceRepository,
	}
}

func (u *useCase) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := u.userRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "email",
				Value:    user.Email,
				Operator: criteria.OperatorEqual,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, apperror.NewInvalidData("User with this email already exists", nil, "user_usecase.go:Create")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *useCase) Get(ctx context.Context, id int) (*domain.UserWithSpaces, error) {
	user, err := u.userRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    id,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperror.NewNotFound("User not found", nil, "user_usecase.go:GetUserWithSpaces")
	}

	spaceIDs, err := u.userSpaceRepository.FindSpacesIDsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	spaces, err := u.spaceRepository.FindByIDs(ctx, spaceIDs)

	if err != nil {
		return nil, err
	}

	return &domain.UserWithSpaces{
		User:   user,
		Spaces: spaces,
	}, nil
}

func (u *useCase) GetSpacesByUser(ctx context.Context, userId int) ([]*domain.Space, error) {
	spaceIDs, err := u.userSpaceRepository.FindSpacesIDsByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return u.spaceRepository.FindByIDs(ctx, spaceIDs)
}

func (u *useCase) Update(ctx context.Context, dto dto.UpdateUserSpacesDTO) error {
	user, err := u.userRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    dto.UserID,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return err
	}

	if user == nil {
		return apperror.NewNotFound("User not found", nil, "user_usecase.go:Update")
	}

	if len(dto.SpaceIDs) == 0 {
		return apperror.NewInvalidData("Space IDs cannot be empty", nil, "user_usecase.go:Update")
	}

	if err := u.userSpaceRepository.Update(ctx, user.ID, dto.SpaceIDs, dto.Action); err != nil {
		return err
	}

	return nil
}
