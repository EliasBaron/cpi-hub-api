package mapper

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/entity"
)

func ToMongoReaction(reaction *domain.Reaction) *entity.Reaction {
	return &entity.Reaction{
		ID:         reaction.ID,
		UserID:     reaction.UserID,
		EntityType: reaction.EntityType,
		EntityID:   reaction.EntityID,
		Liked:      reaction.Liked,
		Disliked:   reaction.Disliked,
	}
}

func ToDomainReaction(reactionEntity *entity.Reaction) *domain.Reaction {
	return &domain.Reaction{
		ID:         reactionEntity.ID,
		UserID:     reactionEntity.UserID,
		EntityType: reactionEntity.EntityType,
		EntityID:   reactionEntity.EntityID,
		Liked:      reactionEntity.Liked,
		Disliked:   reactionEntity.Disliked,
	}
}
