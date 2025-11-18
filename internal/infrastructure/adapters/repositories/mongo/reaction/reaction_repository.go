package reaction

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/mapper"
	"cpi-hub-api/pkg/apperror"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	res, err := collection.InsertOne(ctx, reactionEntity)
	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		reaction.ID = oid.Hex()
	}

	return nil
}

func (r *ReactionRepository) FindReaction(ctx context.Context, criteria *criteria.Criteria) (*domain.Reaction, error) {
	mongoQuery := mapper.ToMongoDBQuery(criteria)

	var reactionEntity entity.Reaction
	err := r.db.Collection("reactions").FindOne(ctx, mongoQuery).Decode(&reactionEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find reaction: %w", err)
	}

	return mapper.ToDomainReaction(&reactionEntity), nil
}

func (r *ReactionRepository) DeleteReaction(ctx context.Context, reactionID string) error {
	oid, err := primitive.ObjectIDFromHex(reactionID)
	if err != nil {
		return fmt.Errorf("invalid reaction ID: %w", err)
	}

	res, err := r.db.Collection("reactions").DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return fmt.Errorf("failed to delete reaction: %w", err)
	}

	if res.DeletedCount == 0 {
		return apperror.NewNotFound("Reaction not found", nil, "reaction_repository.go:DeleteReaction")
	}

	return nil
}

func (r *ReactionRepository) UpdateReaction(ctx context.Context, reaction *domain.Reaction) error {
	oid, err := primitive.ObjectIDFromHex(reaction.ID)
	if err != nil {
		return fmt.Errorf("invalid reaction ID: %w", err)
	}

	update := bson.M{
		"$set": mapper.ToMongoReaction(reaction),
	}

	_, err = r.db.Collection("reactions").UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return fmt.Errorf("failed to update reaction: %w", err)
	}

	return nil
}

func (r *ReactionRepository) CountReactions(ctx context.Context, criteria *criteria.Criteria) (int, error) {
	mongoQuery := mapper.ToMongoDBQuery(criteria)

	count, err := r.db.Collection("reactions").CountDocuments(ctx, mongoQuery)
	if err != nil {
		return 0, fmt.Errorf("failed to count reactions: %w", err)
	}

	return int(count), nil
}
