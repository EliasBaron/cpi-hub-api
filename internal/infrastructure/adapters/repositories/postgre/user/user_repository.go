package user

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgre/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgre/mapper"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func (u *UserRepository) findUserByField(ctx context.Context, field string, value string) (*domain.User, error) {
	var userEntity entity.UserEntity
	query := "SELECT * FROM users WHERE " + field + " = $1"
	err := u.db.QueryRowContext(ctx, query, value).Scan(
		&userEntity.ID, &userEntity.Name, &userEntity.LastName, &userEntity.Email,
		&userEntity.Password, &userEntity.CreatedAt, &userEntity.UpdatedAt, &userEntity.Image,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToDomainUser(&userEntity), nil
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return u.findUserByField(ctx, "email", email)
}
func (u *UserRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	return u.findUserByField(ctx, "id", id)
}

func (u *UserRepository) Create(ctx context.Context, user *domain.User) error {
	var userEntity = *mapper.ToPostgreUser(user)
	err := u.db.QueryRowContext(ctx,
		"INSERT INTO users (name, last_name, email, password, created_at, updated_at, image) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		userEntity.Name, userEntity.LastName, userEntity.Email, userEntity.Password, userEntity.CreatedAt, userEntity.UpdatedAt, userEntity.Image,
	).Scan(&userEntity.ID)
	if err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}
	user.ID = userEntity.ID
	return nil
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) AddSpaceToUser(ctx context.Context, userId string, spaceId string) error {
	_, err := u.db.ExecContext(ctx,
		"INSERT INTO user_space (user_id, space_id) VALUES ($1, $2)",
		userId, spaceId,
	)
	if err != nil {
		return fmt.Errorf("error al agregar espacio a usuario: %w", err)
	}
	return nil
}
