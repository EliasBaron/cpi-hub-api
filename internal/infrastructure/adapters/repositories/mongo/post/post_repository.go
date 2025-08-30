package post

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/mapper"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository struct {
	db *mongo.Database
}

func NewPostRepository(db *mongo.Database) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	postEntity := mapper.ToMongoDatabasePost(post)

	collection := r.db.Collection("posts")
	_, err := collection.InsertOne(ctx, postEntity)

	if err != nil {
		return fmt.Errorf("failed to insert post: %w", err)
	}

	return nil
}
