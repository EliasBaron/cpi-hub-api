package reaction

import (
	"context"
	"cpi-hub-api/internal/core/domain"
)

type ReactionUseCase interface {
	AddReaction(ctx context.Context, reaction *domain.Reaction) (*domain.Reaction, error)
	// RemoveReaction(ctx context.Context, reaction *domain.Reaction) error
	// GetReactions(ctx context.Context, criteria *criteria.Criteria) ([]*domain.Reaction, error)
}

type reactionUsecase struct {
	reactionRepo domain.ReactionRepository
}

func NewReactionUsecase(reactionRepo domain.ReactionRepository) ReactionUseCase {
	return &reactionUsecase{
		reactionRepo: reactionRepo,
	}
}

func (u *reactionUsecase) AddReaction(ctx context.Context, reaction *domain.Reaction) (*domain.Reaction, error) {
	reaction, err := u.reactionRepo.AddReaction(ctx, reaction)

	if err != nil {
		return nil, err
	}

	return reaction, nil
}
