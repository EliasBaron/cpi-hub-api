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
		"INSERT INTO comments (post_id, content, created_by, created_at, updated_at, parent_comment_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		commentEntity.PostID, commentEntity.Content, commentEntity.CreatedBy, commentEntity.CreatedAt, commentEntity.UpdatedAt, commentEntity.ParentID,
	).Scan(&commentEntity.ID); err != nil {
		return err
	}

	comment.ID = commentEntity.ID
	return nil
}

func (c *CommentRepository) Find(ctx context.Context, criteria *criteria.Criteria) (*domain.CommentWithInfo, error) {
	query, params := c.buildQueryWithAliases(criteria)
	query += " LIMIT 1"

	comments, err := c.findWithInfoByField(ctx, query, params)
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, nil
	}
	return comments[0], nil
}

func (c *CommentRepository) FindAll(ctx context.Context, criteria *criteria.Criteria) ([]*domain.CommentWithInfo, error) {
	query, params := c.buildQueryWithAliases(criteria)
	return c.findWithInfoByField(ctx, query, params)
}

func (c *CommentRepository) FindWithSpace(ctx context.Context, criteria *criteria.Criteria) ([]*domain.CommentWithInfo, error) {
	query, params := mapper.ToPostgreSQLQuery(criteria)
	return c.findWithSpaceByField(ctx, query, params)
}

func (c *CommentRepository) findWithSpaceByField(ctx context.Context, whereClause string, params []interface{}) ([]*domain.CommentWithInfo, error) {
	var commentsWithInfo []*domain.CommentWithInfo
	query := `
		SELECT 
			c.id, c.post_id, c.content, c.created_by, c.created_at, c.parent_comment_id,
			u.id, u.name, u.last_name, u.email, u.image, u.created_at,
			s.id, s.name, s.description, s.created_at
		FROM comments c
		INNER JOIN posts p ON c.post_id = p.id
		INNER JOIN spaces s ON p.space_id = s.id
		INNER JOIN users u ON c.created_by = u.id
	` + " " + whereClause

	rows, err := c.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var commentEntity entity.CommentEntity
		var userEntity entity.UserEntity
		var spaceEntity entity.SpaceEntity

		if err := rows.Scan(
			&commentEntity.ID,
			&commentEntity.PostID,
			&commentEntity.Content,
			&commentEntity.CreatedBy,
			&commentEntity.CreatedAt,
			&commentEntity.ParentID,
			&userEntity.ID,
			&userEntity.Name,
			&userEntity.LastName,
			&userEntity.Email,
			&userEntity.Image,
			&userEntity.CreatedAt,
			&spaceEntity.ID,
			&spaceEntity.Name,
			&spaceEntity.Description,
			&spaceEntity.CreatedAt,
		); err != nil {
			return nil, err
		}

		comment := mapper.ToDomainComment(&commentEntity)
		user := mapper.ToDomainUser(&userEntity)
		space := mapper.ToDomainSpace(&spaceEntity)

		commentsWithInfo = append(commentsWithInfo, &domain.CommentWithInfo{
			Comment: comment,
			User:    user,
			Space:   space,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commentsWithInfo, nil
}

func (c *CommentRepository) findWithInfoByField(ctx context.Context, whereClause string, params []interface{}) ([]*domain.CommentWithInfo, error) {
	var commentsWithInfo []*domain.CommentWithInfo
	query := `
		SELECT 
			c.id, c.post_id, c.content, c.created_by, c.created_at, c.parent_comment_id,
			u.id, u.name, u.last_name, u.email, u.image, u.created_at,
			s.id, s.name, s.description, s.created_at
		FROM comments c
		INNER JOIN posts p ON c.post_id = p.id
		INNER JOIN spaces s ON p.space_id = s.id
		INNER JOIN users u ON c.created_by = u.id
	` + " " + whereClause

	rows, err := c.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var commentEntity entity.CommentEntity
		var userEntity entity.UserEntity
		var spaceEntity entity.SpaceEntity

		if err := rows.Scan(
			&commentEntity.ID,
			&commentEntity.PostID,
			&commentEntity.Content,
			&commentEntity.CreatedBy,
			&commentEntity.CreatedAt,
			&commentEntity.ParentID,
			&userEntity.ID,
			&userEntity.Name,
			&userEntity.LastName,
			&userEntity.Email,
			&userEntity.Image,
			&userEntity.CreatedAt,
			&spaceEntity.ID,
			&spaceEntity.Name,
			&spaceEntity.Description,
			&spaceEntity.CreatedAt,
		); err != nil {
			return nil, err
		}

		comment := mapper.ToDomainComment(&commentEntity)
		user := mapper.ToDomainUser(&userEntity)
		space := mapper.ToDomainSpace(&spaceEntity)

		commentsWithInfo = append(commentsWithInfo, &domain.CommentWithInfo{
			Comment: comment,
			User:    user,
			Space:   space,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commentsWithInfo, nil
}

func (c *CommentRepository) buildQueryWithAliases(criteria *criteria.Criteria) (string, []interface{}) {
	return mapper.ToPostgreSQLQueryWithAlias(criteria, "c")
}

func (c *CommentRepository) Count(ctx context.Context, criteria *criteria.Criteria) (int, error) {
	query, params := mapper.ToPostgreSQLCountQuery(criteria)

	countQuery := `
		SELECT COUNT(*)
		FROM comments
	` + " " + query

	var count int
	err := c.db.QueryRowContext(ctx, countQuery, params...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *CommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	var commentEntity = *mapper.ToPostgreComment(comment)

	query := `
		UPDATE comments
		SET content = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := c.db.ExecContext(ctx, query, commentEntity.Content, commentEntity.UpdatedAt, commentEntity.ID)
	return err
}

func (c *CommentRepository) DeleteChildren(ctx context.Context, parentCommentID int) error {
	query := `
		DELETE FROM comments WHERE parent_comment_id = $1
	`
	_, err := c.db.ExecContext(ctx, query, parentCommentID)
	return err
}

func (c *CommentRepository) DeleteFromPost(ctx context.Context, postID int) error {
	query := `
		DELETE FROM comments WHERE post_id = $1
	`
	_, err := c.db.ExecContext(ctx, query, postID)
	return err
}

func (c *CommentRepository) Delete(ctx context.Context, commentID int) error {
	query := `
		DELETE FROM comments
		WHERE id = $1
	`
	_, err := c.db.ExecContext(ctx, query, commentID)
	return err
}
