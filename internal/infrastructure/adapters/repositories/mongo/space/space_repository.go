package space

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
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

func (r *SpaceRepository) Create(ctx context.Context, space *domain.Space) error {
	spaceEntity := mapper.ToMongoDatabaseSpace(space)

	collection := r.db.Collection("spaces")

	_, err := collection.InsertOne(ctx, spaceEntity)

	if err != nil {
		return fmt.Errorf("error al crear el espacio: %w", err)
	}

	return nil
}

func (r *SpaceRepository) Find(ctx context.Context, criteria *criteria.Criteria) (*domain.Space, error) {
	filters := mapper.ToMongoDBQuery(criteria)

	collection := r.db.Collection("spaces")
	cursor, err := collection.Find(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error al buscar el espacio: %w", err)
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

func (r *SpaceRepository) FindByIDs(ctx context.Context, ids []string) ([]*domain.Space, error) {
	if len(ids) == 0 {
		return []*domain.Space{}, nil
	}
	collection := r.db.Collection("spaces")
	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, fmt.Errorf("error al buscar espacios por IDs: %w", err)
	}
	defer cursor.Close(ctx)

	var results []*domain.Space
	for cursor.Next(ctx) {
		var spaceEntity entity.SpaceEntity
		if err := cursor.Decode(&spaceEntity); err != nil {
			return nil, fmt.Errorf("error al decodificar espacio: %w", err)
		}
		results = append(results, mapper.ToDomainSpace(&spaceEntity))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
