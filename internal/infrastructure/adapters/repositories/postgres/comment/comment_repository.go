package comment

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
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

func (c *CommentRepository) Find(ctx context.Context, criteria *criteria.Criteria) (*domain.CommentWithUser, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CommentRepository) FindAllByPostID(ctx context.Context, postID int) ([]*domain.CommentWithUser, error) {
	//TODO implement me
	panic("implement me")
}
