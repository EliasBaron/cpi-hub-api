package user

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/mapper"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	userEntity := mapper.ToMongoDatabaseUser(user)
	collection := r.db.Collection("users")
	_, err := collection.InsertOne(ctx, userEntity)
	if err != nil {
		return fmt.Errorf("error al crear el usuario: %w", err)
	}
	return nil
}
