package user

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/pkg/apperror"
	"strings"
	"time"
)

type UserUseCase interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Get(ctx context.Context, id string) (*domain.UserWithSpaces, error)
	AddSpaceToUser(ctx context.Context, userId string, spaceId string) error
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
	existingUser, err := u.userRepository.FindByEmail(ctx, user.Email)
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

func (u *useCase) Get(ctx context.Context, id string) (*domain.UserWithSpaces, error) {
	user, err := u.userRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NewNotFound("User not found", nil, "user_usecase.go:GetUserWithSpaces")
	}

	spaceIDs, err := u.userSpaceRepository.FindSpaceIDsByUser(ctx, user.ID)
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

func (u *useCase) AddSpaceToUser(ctx context.Context, userId string, spaceId string) error {
	user, err := u.userRepository.FindById(ctx, userId)
	if err != nil {
		return err
	}
	if user == nil {
		return apperror.NewNotFound("User not found", nil, "user_usecase.go:AddSpaceToUser")
	}

	space, err := u.spaceRepository.FindById(ctx, spaceId)
	if err != nil {
		return err
	}
	if space == nil {
		return apperror.NewNotFound("Space not found", nil, "user_usecase.go:AddSpaceToUser")
	}

	err = u.userSpaceRepository.AddUserToSpace(ctx, userId, spaceId)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return apperror.NewInvalidData("User already subscribed to this space", nil, "user_usecase.go:AddSpaceToUser")
		}
		return err
	}

	return nil
}
