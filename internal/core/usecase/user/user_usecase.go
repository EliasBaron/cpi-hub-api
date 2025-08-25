package user

import "cpi-hub-api/internal/core/domain"

type UseCase interface {
	Create(user *domain.User) error
}

type useCase struct {
}

func NewUserUsecase() UseCase {
	return &useCase{}
}

func (u *useCase) Create(user *domain.User) error {
	return nil
}
