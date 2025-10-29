package reaction

import (
	"context"
	"cpi-hub-api/internal/core/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type ReactionRepository struct {
	db *mongo.Database
}

func NewReactionRepository(db *mongo.Database) *ReactionRepository {
	return &ReactionRepository{
		db: db,
	}
}

func (r *ReactionRepository) AddReaction(ctx context.Context, reaction *domain.Reaction) error {
	collection := r.db.Collection("reactions")
	_, err := collection.InsertOne(ctx, reaction)
	return err
}
