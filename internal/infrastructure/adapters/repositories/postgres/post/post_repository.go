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
	posts, err := p.executeQuery(ctx, query+" LIMIT 1", params)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, nil
	}
	return posts[0], nil
}

func (p *PostRepository) FindAll(ctx context.Context, criteria *criteria.Criteria) ([]*domain.Post, error) {
	query, params := mapper.ToPostgreSQLQuery(criteria)
	return p.executeQuery(ctx, query, params)
}

func (p *PostRepository) executeQuery(ctx context.Context, whereClause string, params []interface{}) ([]*domain.Post, error) {
	var posts []*domain.Post

	query := `
		SELECT id, title, content, created_by, created_at, updated_by, updated_at, space_id
		FROM posts
	` + " " + whereClause + " ORDER BY created_at DESC"

	rows, err := p.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postEntity entity.PostEntity
		if err := p.scanPostEntity(rows, &postEntity); err != nil {
			return nil, err
		}
		posts = append(posts, mapper.ToDomainPost(&postEntity))
	}

	return posts, rows.Err()
}

func (p *PostRepository) scanPostEntity(scanner interface{ Scan(...interface{}) error }, postEntity *entity.PostEntity) error {
	return scanner.Scan(
		&postEntity.ID,
		&postEntity.Title,
		&postEntity.Content,
		&postEntity.CreatedBy,
		&postEntity.CreatedAt,
		&postEntity.UpdatedBy,
		&postEntity.UpdatedAt,
		&postEntity.SpaceID,
	)
}

func (p *PostRepository) SearchByTitleOrContent(ctx context.Context, query string) ([]*domain.Post, error) {
	sqlQuery := `
        SELECT id, title, content, created_by, created_at, updated_by, updated_at, space_id
        FROM posts
        WHERE title ILIKE '%' || $1 || '%' OR content ILIKE '%' || $1 || '%'
    `

	rows, err := p.db.QueryContext(ctx, sqlQuery, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var postEntity entity.PostEntity
		if err := p.scanPostEntity(rows, &postEntity); err != nil {
			return nil, err
		}
		posts = append(posts, mapper.ToDomainPost(&postEntity))
	}

	return posts, rows.Err()
}
