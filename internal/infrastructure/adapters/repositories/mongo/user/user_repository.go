package user

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/mapper"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) findByField(ctx context.Context, field string, value interface{}) (*domain.User, error) {
	collection := r.db.Collection("users")
	cursor, err := collection.Find(ctx, bson.D{{Key: field, Value: value}})
	if err != nil {
		return nil, fmt.Errorf("error al buscar el usuario por %s: %w", field, err)
	}
	defer cursor.Close(ctx)

	var userEntity entity.UserEntity
	if cursor.Next(ctx) {
		if err := cursor.Decode(&userEntity); err != nil {
			return nil, fmt.Errorf("error al decodificar el usuario: %w", err)
		}
		return mapper.ToDomainUser(&userEntity), nil
	}
	return nil, mongo.ErrNoDocuments
}

func (r *UserRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	return r.findByField(ctx, "_id", id)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.findByField(ctx, "email", email)
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {

	existingUser, err := r.FindByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("error checking for existing user: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("error: user with email %s already exists", user.Email)
	}

	userEntity := mapper.ToMongoDatabaseUser(user)
	collection := r.db.Collection("users")
	_, err = collection.InsertOne(ctx, userEntity)
	if err != nil {
		return fmt.Errorf("error al crear el usuario: %w", err)
	}
	return nil
}
