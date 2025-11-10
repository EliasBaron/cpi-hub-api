package dto

import (
	"cpi-hub-api/internal/core/domain"
)

type NewReaction struct {
	UserID     int    `json:"user_id" binding:"required"`
	EntityType string `json:"entity_type" binding:"required"`
	EntityID   int    `json:"entity_id" binding:"required"`
	Action     string `json:"action" binding:"required"`
}

type ReactionDTO struct {
	ID         string `json:"id"`
	UserID     int    `json:"user_id"`
	EntityType string `json:"entity_type"`
	EntityID   int    `json:"entity_id"`
	Action     string `json:"action"`
}

type GetLikesCountDTO struct {
	EntityType *string `json:"entity_type"`
	EntityID   *int    `json:"entity_id"`
	UserID     *int    `json:"user_id"`
}

type LikesCountDTO struct {
	UserID        *int    `json:"user_id,omitempty"`
	EntityType    *string `json:"entity_type,omitempty"`
	EntityID      *int    `json:"entity_id,omitempty"`
	LikesCount    int     `json:"likes_count"`
	DislikesCount int     `json:"dislikes_count"`
}

type EntityDataDTO struct {
	EntityType domain.EntityType `json:"entity_type"`
	EntityID   int               `json:"entity_id"`
}

type EntitiesDataDTO struct {
	Entities []EntityDataDTO `json:"entities" binding:"required"`
}

type UserLikeDTO struct {
	EntityType string `json:"entity_type"`
	EntityID   int    `json:"entity_id"`
	Liked      bool   `json:"liked"`
	Disliked   bool   `json:"disliked"`
}

func ToReactionDTO(reaction domain.Reaction) ReactionDTO {
	return ReactionDTO{
		ID:         reaction.ID,
		UserID:     reaction.UserID,
		EntityType: string(reaction.EntityType),
		EntityID:   reaction.EntityID,
		Action:     string(reaction.Action),
	}
}

func (r *ReactionDTO) ToDomain() domain.Reaction {
	return domain.Reaction{
		ID:         r.ID,
		UserID:     r.UserID,
		EntityType: domain.EntityType(r.EntityType),
		EntityID:   r.EntityID,
		Action:     domain.ActionType(r.Action),
	}
}

func (n *NewReaction) ToDomain() *domain.Reaction {
	return &domain.Reaction{
		UserID:     n.UserID,
		EntityType: domain.EntityType(n.EntityType),
		EntityID:   n.EntityID,
		Action:     domain.ActionType(n.Action),
	}
}
