package user

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/mapper"
	"cpi-hub-api/pkg/apperror"
	"database/sql"
	"fmt"
)

type UserSpaceRepository struct {
	db *sql.DB
}

func NewUserSpaceRepository(db *sql.DB) *UserSpaceRepository {
	return &UserSpaceRepository{db: db}
}

func (r *UserSpaceRepository) findIDsByField(ctx context.Context, selectField, whereField string, whereValue int) ([]int, error) {
	query := fmt.Sprintf(`SELECT %s FROM user_spaces WHERE %s = $1`, selectField, whereField)

	rows, err := r.db.QueryContext(ctx, query, whereValue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if ids == nil {
		return []int{}, nil
	}

	return ids, nil
}

func (r *UserSpaceRepository) FindSpacesIDsByUserID(ctx context.Context, userID int) ([]int, error) {
	spaceIDs, err := r.findIDsByField(ctx, "space_id", "user_id", userID)
	if err != nil {
		return nil, err
	}

	return spaceIDs, nil
}

func (r *UserSpaceRepository) FindUserIDsBySpaceID(ctx context.Context, spaceID int) ([]int, error) {
	return r.findIDsByField(ctx, "user_id", "space_id", spaceID)
}

func (u *UserSpaceRepository) Update(ctx context.Context, userId int, spaceIDs []int, action string) error {
	switch action {
	case domain.AddUserToSpace:
		for _, spaceID := range spaceIDs {
			_, err := u.db.ExecContext(ctx, "INSERT INTO user_spaces (user_id, space_id) VALUES ($1, $2)", userId, spaceID)
			if err != nil {
				return err
			}
		}
	case domain.RemoveUserFromSpace:
		for _, spaceID := range spaceIDs {
			_, err := u.db.ExecContext(ctx, "DELETE FROM user_spaces WHERE user_id = $1 AND space_id = $2", userId, spaceID)
			if err != nil {
				return err
			}
		}
	default:
		return apperror.NewInvalidData("invalid action", nil, "user_space_repository.go:Update")
	}

	return nil
}

func (u *UserSpaceRepository) Exists(ctx context.Context, userId int, spaceId int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user_spaces WHERE user_id = $1 AND space_id = $2)`
	err := u.db.QueryRowContext(ctx, query, userId, spaceId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p *UserSpaceRepository) Count(ctx context.Context, criteria *criteria.Criteria) (int, error) {
	whereClause, params := mapper.ToPostgreSQLQuery(criteria)

	query := "SELECT COUNT(*) FROM user_spaces"
	if whereClause != "" {
		query += whereClause
	}

	var count int
	err := p.db.QueryRowContext(ctx, query, params...).Scan(&count)
	return count, err
}
