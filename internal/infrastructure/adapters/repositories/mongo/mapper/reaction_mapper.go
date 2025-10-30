package mapper

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToMongoReaction(reaction *domain.Reaction) *entity.Reaction {
	var oid primitive.ObjectID
	if reaction.ID != "" {
		if parsed, err := primitive.ObjectIDFromHex(reaction.ID); err == nil {
			oid = parsed
		}
	}
	return &entity.Reaction{
		ID:         oid,
		UserID:     reaction.UserID,
		EntityType: string(reaction.EntityType),
		EntityID:   reaction.EntityID,
		Liked:      reaction.Liked,
		Disliked:   reaction.Disliked,
	}
}

func ToDomainReaction(reactionEntity *entity.Reaction) *domain.Reaction {
	var idStr string
	if reactionEntity.ID != primitive.NilObjectID {
		idStr = reactionEntity.ID.Hex()
	}
	return &domain.Reaction{
		ID:         idStr,
		UserID:     reactionEntity.UserID,
		EntityType: domain.EntityType(reactionEntity.EntityType),
		EntityID:   reactionEntity.EntityID,
		Liked:      reactionEntity.Liked,
		Disliked:   reactionEntity.Disliked,
	}
}
