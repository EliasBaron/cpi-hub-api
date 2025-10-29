package reaction

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/mapper"
	"fmt"

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

	reactionEntity := mapper.ToMongoReaction(reaction)

	collection := r.db.Collection("reactions")

	_, err := collection.InsertOne(ctx, reactionEntity)

	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	return nil
}
