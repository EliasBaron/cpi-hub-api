package space

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/mapper"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SpaceRepository struct {
	db *mongo.Database
}

func NewSpaceRepository(db *mongo.Database) *SpaceRepository {
	return &SpaceRepository{db: db}
}

func (r *SpaceRepository) findByField(ctx context.Context, field string, value interface{}) (*domain.Space, error) {
	collection := r.db.Collection("spaces")
	cursor, err := collection.Find(ctx, bson.D{{Key: field, Value: value}})
	if err != nil {
		return nil, fmt.Errorf("error al buscar el espacio por %s: %w", field, err)
	}
	defer cursor.Close(ctx)

	var spaceEntity entity.SpaceEntity
	if cursor.Next(ctx) {
		if err := cursor.Decode(&spaceEntity); err != nil {
			return nil, fmt.Errorf("error al decodificar el espacio: %w", err)
		}
		return mapper.ToDomainSpace(&spaceEntity), nil
	}
	return nil, nil
}

func (r *SpaceRepository) FindById(ctx context.Context, id string) (*domain.Space, error) {
	return r.findByField(ctx, "_id", id)
}

func (r *SpaceRepository) FindByName(ctx context.Context, name string) (*domain.Space, error) {
	return r.findByField(ctx, "name", name)
}

func (r *SpaceRepository) Create(ctx context.Context, space *domain.Space) error {
	spaceEntity := mapper.ToMongoDatabaseSpace(space)
	collection := r.db.Collection("spaces")
	_, err := collection.InsertOne(ctx, spaceEntity)
	if err != nil {
		return fmt.Errorf("error al crear el espacio: %w", err)
	}
	return nil
}
