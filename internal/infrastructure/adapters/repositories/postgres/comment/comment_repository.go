package comment

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/mapper"
	"database/sql"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (c *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	var commentEntity = *mapper.ToPostgreComment(comment)

	if err := c.db.QueryRowContext(ctx,
		"INSERT INTO comments (post_id, content, created_by, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		commentEntity.PostID, commentEntity.Content, commentEntity.CreatedBy, commentEntity.CreatedAt,
	).Scan(&commentEntity.ID); err != nil {
		return err
	}

	comment.ID = commentEntity.ID
	return nil
}

func (c *CommentRepository) FindAll(ctx context.Context, criteria *criteria.Criteria) ([]*domain.Comment, error) {
	query, params := mapper.ToPostgreSQLQuery(criteria)
	return c.findAllByField(ctx, query, params)
}

func (c *CommentRepository) findAllByField(ctx context.Context, whereClause string, params []interface{}) ([]*domain.Comment, error) {
	var comments []*domain.Comment
	query := `
		SELECT id, post_id, content, created_by, created_at
		FROM comments
	` + " " + whereClause

	rows, err := c.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var commentEntity entity.CommentEntity
		if err := rows.Scan(
			&commentEntity.ID,
			&commentEntity.PostID,
			&commentEntity.Content,
			&commentEntity.CreatedBy,
			&commentEntity.CreatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, mapper.ToDomainComment(&commentEntity))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentRepository) Find(ctx context.Context, criteria *criteria.Criteria) (*domain.Comment, error) {
	comments, err := c.FindAll(ctx, criteria)
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, nil
	}
	return comments[0], nil
}
