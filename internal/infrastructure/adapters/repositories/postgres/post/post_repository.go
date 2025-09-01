package post

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/mapper"
	"database/sql"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (p *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	var postEntity = *mapper.ToPostgresPost(post)

	if err := p.db.QueryRowContext(ctx,
		"INSERT INTO posts (title, content, created_at, updated_at, created_by, updated_by, space_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		postEntity.Title, postEntity.Content, postEntity.CreatedAt, postEntity.UpdatedAt, postEntity.CreatedBy, postEntity.UpdatedBy, postEntity.SpaceID).Scan(&postEntity.ID); err != nil {
		return err
	}

	post.ID = postEntity.ID
	return nil
}

func (p *PostRepository) Find(ctx context.Context, criteria *criteria.Criteria) (*domain.Post, error) {
	query, params := mapper.ToPostgreSQLQuery(criteria)

	return p.findPostByField(ctx, query, params)
}

func (p *PostRepository) findPostByField(ctx context.Context, whereClause string, params []interface{}) (*domain.Post, error) {
	var postEntity entity.PostEntity

	query := `
        SELECT id, title, content, created_by, created_at, updated_by, updated_at, space_id
        FROM posts
    ` + " " + whereClause + " LIMIT 1"

	err := p.db.QueryRowContext(ctx, query, params...).Scan(
		&postEntity.ID,
		&postEntity.Title,
		&postEntity.Content,
		&postEntity.CreatedBy,
		&postEntity.CreatedAt,
		&postEntity.UpdatedBy,
		&postEntity.UpdatedAt,
		&postEntity.SpaceID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToDomainPost(&postEntity), nil
}
