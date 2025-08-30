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
	FindByIDs(ctx context.Context, ids []string) ([]*Space, error)
}

type UserSpaceRepository interface {
	AddUserToSpace(ctx context.Context, userId int, spaceId string) error
	FindSpaceIDsByUser(ctx context.Context, userID int) ([]string, error)
	Exists(ctx context.Context, userId int, spaceId string) (bool, error)
}
