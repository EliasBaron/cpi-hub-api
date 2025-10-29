package dto

import (
	"cpi-hub-api/internal/core/domain"
)

type NewReaction struct {
	UserID     int    `json:"user_id" binding:"required"`
	EntityType string `json:"entity_type" binding:"required"`
	EntityID   int    `json:"entity_id" binding:"required"`
	Liked      bool   `json:"liked"`
	Disliked   bool   `json:"disliked"`
}

type ReactionDTO struct {
	ID         string `json:"id"`
	UserID     int    `json:"user_id"`
	EntityType string `json:"entity_type"`
	EntityID   int    `json:"entity_id"`
	Liked      bool   `json:"liked"`
	Disliked   bool   `json:"disliked"`
}

func ToReactionDTO(reaction domain.Reaction) ReactionDTO {
	return ReactionDTO{
		ID:         reaction.ID,
		UserID:     reaction.UserID,
		EntityType: reaction.EntityType,
		EntityID:   reaction.EntityID,
		Liked:      reaction.Liked,
		Disliked:   reaction.Disliked,
	}
}

func (r *ReactionDTO) ToDomain() domain.Reaction {
	return domain.Reaction{
		ID:         r.ID,
		UserID:     r.UserID,
		EntityType: r.EntityType,
		EntityID:   r.EntityID,
		Liked:      r.Liked,
		Disliked:   r.Disliked,
	}
}

func (n *NewReaction) ToDomain() *domain.Reaction {
	return &domain.Reaction{
		UserID:     n.UserID,
		EntityType: n.EntityType,
		EntityID:   n.EntityID,
		Liked:      n.Liked,
		Disliked:   n.Disliked,
	}
}
