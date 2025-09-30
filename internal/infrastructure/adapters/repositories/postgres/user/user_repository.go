package user

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/mapper"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Create(ctx context.Context, user *domain.User) error {
	var userEntity = *mapper.ToPostgreUser(user)

	if err := u.db.QueryRowContext(ctx,
		"INSERT INTO users (name, last_name, email, password, created_at, updated_at, image) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		userEntity.Name, userEntity.LastName, userEntity.Email, userEntity.Password, userEntity.CreatedAt, userEntity.UpdatedAt, userEntity.Image,
	).Scan(&userEntity.ID); err != nil {
		return err
	}

	user.ID = userEntity.ID
	return nil
}

func (u *UserRepository) Find(ctx context.Context, criteria *criteria.Criteria) (*domain.User, error) {
	query, params := mapper.ToPostgreSQLQuery(criteria)

	return u.findUserByField(ctx, query, params)
}

func (u *UserRepository) findUserByField(ctx context.Context, whereClause string, params []interface{}) (*domain.User, error) {
	var userEntity entity.UserEntity

	query := `
		SELECT id, name, last_name, email, password, created_at, updated_at, image
		FROM users
	` + " " + whereClause + " LIMIT 1"

	err := u.db.QueryRowContext(ctx, query, params...).Scan(
		&userEntity.ID,
		&userEntity.Name,
		&userEntity.LastName,
		&userEntity.Email,
		&userEntity.Password,
		&userEntity.CreatedAt,
		&userEntity.UpdatedAt,
		&userEntity.Image,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToDomainUser(&userEntity), nil
}

func (u *UserRepository) Search(ctx context.Context, criteria *criteria.Criteria) ([]*domain.User, error) {
	query, params := mapper.ToPostgreSQLQuery(criteria)

	return u.findUsersByField(ctx, query, params)
}

func (u *UserRepository) Count(ctx context.Context, criteria *criteria.Criteria) (int, error) {
	query, params := mapper.ToPostgreSQLQuery(criteria)

	var count int
	countQuery := "SELECT COUNT(*) FROM users " + query

	err := u.db.QueryRowContext(ctx, countQuery, params...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (u *UserRepository) findUsersByField(ctx context.Context, whereClause string, params []interface{}) ([]*domain.User, error) {
	query := `
		SELECT id, name, last_name, email, password, created_at, updated_at, image
		FROM users
	` + " " + whereClause

	rows, err := u.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var userEntity entity.UserEntity
		err := rows.Scan(
			&userEntity.ID,
			&userEntity.Name,
			&userEntity.LastName,
			&userEntity.Email,
			&userEntity.Password,
			&userEntity.CreatedAt,
			&userEntity.UpdatedAt,
			&userEntity.Image,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, mapper.ToDomainUser(&userEntity))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
