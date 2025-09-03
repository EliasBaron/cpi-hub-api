package domain

import (
	"context"
	"cpi-hub-api/internal/core/domain/criteria"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Find(ctx context.Context, criteria *criteria.Criteria) (*User, error)
}

type SpaceRepository interface {
	Create(ctx context.Context, space *Space) error
	Find(ctx context.Context, criteria *criteria.Criteria) (*Space, error)
	FindByIDs(ctx context.Context, ids []int) ([]*Space, error)
}

type UserSpaceRepository interface {
	AddUserToSpace(ctx context.Context, userId int, spaceId int) error
	FindSpaceIDsByUser(ctx context.Context, userID int) ([]int, error)
	Exists(ctx context.Context, userId int, spaceId int) (bool, error)
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	Find(ctx context.Context, criteria *criteria.Criteria) (*Post, error)
	FindAll(ctx context.Context, criteria *criteria.Criteria) ([]*Post, error)
	SearchByTitleOrContent(ctx context.Context, query string) ([]*Post, error)
}

type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	FindAll(ctx context.Context, c *criteria.Criteria) ([]*Comment, error)
}
