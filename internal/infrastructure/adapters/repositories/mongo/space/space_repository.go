package space

import (
	"context"
	"cpi-hub-api/internal/core/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type SpaceRepository struct {
	db *mongo.Database
}

func NewSpaceRepository(db *mongo.Database) *SpaceRepository {
	return &SpaceRepository{db: db}
}

func (r *SpaceRepository) Create(ctx context.Context, space *domain.Space) error {
	return nil
}

func (r *SpaceRepository) FindById(ctx context.Context, id string) (*domain.Space, error) {
	return nil, nil
}

func (r *SpaceRepository) FindByName(ctx context.Context, name string) (*domain.Space, error) {
	return nil, nil
}
