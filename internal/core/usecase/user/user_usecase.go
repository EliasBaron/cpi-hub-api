package user

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
	"time"
)

type UseCase interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
}

type useCase struct {
	userRepository domain.UserRepository
}

func NewUserUsecase(userRepository domain.UserRepository) UseCase {
	return &useCase{
		userRepository: userRepository,
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

	user.ID = helpers.NewULID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *useCase) Get(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.userRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperror.NewNotFound("User not found", nil, "user_usecase.go:Get")
	}

	return user, nil
}
