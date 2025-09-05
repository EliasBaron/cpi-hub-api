package space

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/mapper"
	"database/sql"
)

type SpaceRepository struct {
	db *sql.DB
}

func NewSpaceRepository(db *sql.DB) *SpaceRepository {
	return &SpaceRepository{db: db}
}

func (u *SpaceRepository) Create(ctx context.Context, space *domain.Space) error {
	var spaceEntity = *mapper.ToPostgresSpace(space)

	if err := u.db.QueryRowContext(ctx,
		"INSERT INTO spaces (name, description, created_by, created_at, updated_by, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		spaceEntity.Name, spaceEntity.Description, spaceEntity.CreatedBy, spaceEntity.CreatedAt, spaceEntity.UpdatedBy, spaceEntity.UpdatedAt).Scan(&spaceEntity.ID); err != nil {
		return err
	}

	space.ID = spaceEntity.ID
	return nil
}

func (u *SpaceRepository) Find(ctx context.Context, criteria *criteria.Criteria) (*domain.Space, error) {
	query, params := mapper.ToPostgreSQLQuery(criteria)

	return u.findSpaceByField(ctx, query, params)
}

func (u *SpaceRepository) findSpaceByField(ctx context.Context, whereClause string, params []interface{}) (*domain.Space, error) {
	var spaceEntity entity.SpaceEntity

	query := `
		SELECT id, name, description, created_by, created_at, updated_by, updated_at
		FROM spaces
	` + " " + whereClause + " LIMIT 1"

	if err := u.db.QueryRowContext(ctx, query, params...).Scan(
		&spaceEntity.ID,
		&spaceEntity.Name,
		&spaceEntity.Description,
		&spaceEntity.CreatedBy,
		&spaceEntity.CreatedAt,
		&spaceEntity.UpdatedBy,
		&spaceEntity.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToDomainSpace(&spaceEntity), nil
}

func (u *SpaceRepository) FindByIDs(ctx context.Context, ids []int) ([]*domain.Space, error) {
	if len(ids) == 0 {
		return []*domain.Space{}, nil
	}

	idsInterface := make([]interface{}, len(ids))
	for i, id := range ids {
		idsInterface[i] = id
	}

	spaces, err := u.FindAll(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    idsInterface,
				Operator: criteria.OperatorIn,
			},
		},
	})

	return spaces, err
}

func (u *SpaceRepository) FindAll(ctx context.Context, criteria *criteria.Criteria) ([]*domain.Space, error) {
	whereClause, params := mapper.ToPostgreSQLQuery(criteria)

	query := `
        SELECT id, name, description, created_by, created_at, updated_by, updated_at
        FROM spaces
    ` + " " + whereClause

	rows, err := u.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spaces []*domain.Space
	for rows.Next() {
		var spaceEntity entity.SpaceEntity
		if err := rows.Scan(
			&spaceEntity.ID,
			&spaceEntity.Name,
			&spaceEntity.Description,
			&spaceEntity.CreatedBy,
			&spaceEntity.CreatedAt,
			&spaceEntity.UpdatedBy,
			&spaceEntity.UpdatedAt,
		); err != nil {
			return nil, err
		}
		spaces = append(spaces, mapper.ToDomainSpace(&spaceEntity))
	}

	return spaces, rows.Err()
}

func (u *SpaceRepository) Update(ctx context.Context, space *domain.Space) error {
	spaceEntity := mapper.ToPostgresSpace(space)

	_, err := u.db.ExecContext(ctx,
		"UPDATE spaces SET name=$1, description=$2, updated_by=$3, updated_at=$4 WHERE id=$5",
		spaceEntity.Name, spaceEntity.Description, spaceEntity.UpdatedBy, spaceEntity.UpdatedAt, spaceEntity.ID)

	return err
}
