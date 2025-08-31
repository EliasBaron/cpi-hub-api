package user

import (
	"context"
	"database/sql"
	"fmt"
)

type UserSpaceRepository struct {
	db *sql.DB
}

func NewUserSpaceRepository(db *sql.DB) *UserSpaceRepository {
	return &UserSpaceRepository{db: db}
}

// Devuelve todos los IDs de espacios asociados a un usuario
func (r *UserSpaceRepository) FindSpaceIDsByUser(ctx context.Context, userID int) ([]int, error) {
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

func (u *UserSpaceRepository) AddUserToSpace(ctx context.Context, userId int, spaceId int) error {
	_, err := u.db.ExecContext(ctx,
		"INSERT INTO user_spaces (user_id, space_id) VALUES ($1, $2)",
		userId, spaceId,
	)
	if err != nil {
		return fmt.Errorf("error al agregar espacio a usuario: %w", err)
	}
	return nil
}

func (u *UserSpaceRepository) Exists(ctx context.Context, userId int, spaceId int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user_spaces WHERE user_id = $1 AND space_id = $2)`
	err := u.db.QueryRowContext(ctx, query, userId, spaceId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error al verificar existencia de espacio: %w", err)
	}
	return exists, nil
}
