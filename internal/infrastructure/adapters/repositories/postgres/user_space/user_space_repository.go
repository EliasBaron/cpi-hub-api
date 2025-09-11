package user

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/pkg/apperror"
	"database/sql"
)

type UserSpaceRepository struct {
	db *sql.DB
}

func NewUserSpaceRepository(db *sql.DB) *UserSpaceRepository {
	return &UserSpaceRepository{db: db}
}

func (r *UserSpaceRepository) FindSpacesIDsByUserID(ctx context.Context, userID int) ([]int, error) {
	query := `SELECT space_id FROM user_spaces WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spaceIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		spaceIDs = append(spaceIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if spaceIDs == nil {
		return []int{}, nil
	}

	return spaceIDs, nil

}

func (u *UserSpaceRepository) EditUserSpaces(ctx context.Context, userId int, spaceIDs []int, action string) error {
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
		return apperror.NewInvalidData("invalid action", nil, "user_space_repository.go:EditUserSpaces")
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
