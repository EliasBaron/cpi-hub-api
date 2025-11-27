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

func (r *ReactionRepository) GetTopReactionEntities(ctx context.Context, crit *criteria.Criteria, groupBy string) ([]*domain.TopReactionEntity, int, error) {
	if groupBy != "entity_id" && groupBy != "user_id" {
		return nil, 0, fmt.Errorf("invalid groupBy field: %s (must be 'entity_id' or 'user_id')", groupBy)
	}

	matchQuery := mapper.ToMongoDBQuery(crit)
	matchStage := bson.D{{Key: "$match", Value: matchQuery}}

	// Group by the specified field
	groupStage := bson.D{{Key: "$group", Value: bson.M{
		"_id":   "$" + groupBy,
		"count": bson.M{"$sum": 1},
	}}}

	// Sort by count descending
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"count": -1}}}

	// Build pagination stages from criteria
	var paginationStages []bson.D
	if crit.Pagination.Page > 0 && crit.Pagination.PageSize > 0 {
		skip := (crit.Pagination.Page - 1) * crit.Pagination.PageSize
		paginationStages = []bson.D{
			{{Key: "$skip", Value: int64(skip)}},
			{{Key: "$limit", Value: int64(crit.Pagination.PageSize)}},
		}
	}

	// Build and execute aggregation pipeline
	pipeline := mongo.Pipeline{matchStage, groupStage, sortStage}
	pipeline = append(pipeline, paginationStages...)

	cursor, err := r.db.Collection("reactions").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to aggregate reactions: %w", err)
	}
	defer cursor.Close(ctx)

	// Parse results
	var results []*domain.TopReactionEntity
	for cursor.Next(ctx) {
		var doc struct {
			ID    int `bson:"_id"`
			Count int `bson:"count"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, fmt.Errorf("failed to decode aggregation result: %w", err)
		}

		topEntity := &domain.TopReactionEntity{
			Count:   doc.Count,
			GroupBy: groupBy,
		}
		if groupBy == "entity_id" {
			topEntity.EntityID = doc.ID
		} else {
			topEntity.UserID = doc.ID
		}
		results = append(results, topEntity)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}

	// Get total count of distinct groups
	totalPipeline := mongo.Pipeline{matchStage, groupStage, bson.D{{Key: "$count", Value: "total"}}}
	totalCursor, err := r.db.Collection("reactions").Aggregate(ctx, totalPipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to aggregate total reactions: %w", err)
	}
	defer totalCursor.Close(ctx)

	var total int
	if totalCursor.Next(ctx) {
		var totalDoc struct {
			Total int `bson:"total"`
		}
		if err := totalCursor.Decode(&totalDoc); err == nil {
			total = totalDoc.Total
		}
	}

	return results, total, nil
}
