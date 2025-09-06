package post

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/mapper"
	"database/sql"
	"fmt"
)

type PostRepository struct {
	db *sql.DB
}

type QueryParams struct {
	WhereClause string
	OrderClause string
	LimitClause string
	Args        []interface{}
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
	whereClause, params := mapper.ToPostgreSQLQuery(criteria)
	queryParams := QueryParams{
		WhereClause: whereClause,
		OrderClause: "ORDER BY created_at DESC",
		LimitClause: "LIMIT 1",
		Args:        params,
	}
	posts, err := p.executeQuery(ctx, queryParams)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, nil
	}
	return posts[0], nil
}

func (p *PostRepository) FindAll(ctx context.Context, criteria *criteria.Criteria) ([]*domain.Post, error) {
	fmt.Println("criteria", criteria)
	whereClause, params := mapper.ToPostgreSQLQuery(criteria)
	queryParams := QueryParams{
		WhereClause: whereClause,
		Args:        params,
	}
	return p.executeQuery(ctx, queryParams)
}

func (p *PostRepository) executeQuery(ctx context.Context, params QueryParams) ([]*domain.Post, error) {
	var posts []*domain.Post

	query := "SELECT id, title, content, created_by, created_at, updated_by, updated_at, space_id FROM posts"

	if params.WhereClause != "" {
		query += params.WhereClause
	}
	if params.OrderClause != "" {
		query += " " + params.OrderClause
	}
	if params.LimitClause != "" {
		query += " " + params.LimitClause
	}

	rows, err := p.db.QueryContext(ctx, query, params.Args...)
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

func (p *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	var postEntity = *mapper.ToPostgresPost(post)

	_, err := p.db.ExecContext(ctx,
		"UPDATE posts SET title=$1, content=$2, updated_at=$3, updated_by=$4, space_id=$5 WHERE id=$6",
		postEntity.Title, postEntity.Content, postEntity.UpdatedAt, postEntity.UpdatedBy, postEntity.SpaceID, postEntity.ID)
	return err
}
