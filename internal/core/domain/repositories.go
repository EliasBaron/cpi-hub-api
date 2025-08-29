package domain

import "context"

type UserRepository interface {
	FindById(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
}

type SpaceRepository interface {
	FindById(ctx context.Context, id string) (*Space, error)
	FindByIDs(ctx context.Context, ids []string) ([]*Space, error)
	Create(ctx context.Context, space *Space) error
	FindByName(ctx context.Context, name string) (*Space, error)
}

type UserSpaceRepository interface {
	FindSpaceIDsByUser(ctx context.Context, userID string) ([]string, error)
	AddUserToSpace(ctx context.Context, userId string, spaceId string) error
	Exists(ctx context.Context, userId string, spaceId string) (bool, error)
}
