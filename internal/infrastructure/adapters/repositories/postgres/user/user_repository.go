package user

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

// Create implements domain.UserRepository.
func (u *UserRepository) Create(ctx context.Context, user *domain.User) error {
	panic("unimplemented")
}

// FindByEmail implements domain.UserRepository.
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	panic("unimplemented")
}

// FindById implements domain.UserRepository.
func (u *UserRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	panic("unimplemented")
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}
