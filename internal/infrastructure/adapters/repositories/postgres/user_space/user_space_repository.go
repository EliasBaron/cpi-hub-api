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
func (r *UserSpaceRepository) FindSpaceIDsByUser(ctx context.Context, userID string) ([]string, error) {
	query := `SELECT space_id FROM user_spaces WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spaceIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		spaceIDs = append(spaceIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return spaceIDs, nil
}

func (u *UserSpaceRepository) AddUserToSpace(ctx context.Context, userId string, spaceId string) error {
	_, err := u.db.ExecContext(ctx,
		"INSERT INTO user_spaces (user_id, space_id) VALUES ($1, $2)",
		userId, spaceId,
	)
	if err != nil {
		return fmt.Errorf("error al agregar espacio a usuario: %w", err)
	}
	return nil
}
