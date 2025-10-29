package reaction

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	pghelpers "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/helpers"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
)

type ReactionUseCase interface {
	AddReaction(ctx context.Context, reaction *domain.Reaction) (*domain.Reaction, error)
	// RemoveReaction(ctx context.Context, reaction *domain.Reaction) error
	// GetReactions(ctx context.Context, criteria *criteria.Criteria) ([]*domain.Reaction, error)
}

type reactionUsecase struct {
	reactionRepo domain.ReactionRepository
	userRepo     domain.UserRepository
	postRepo     domain.PostRepository
	commentRepo  domain.CommentRepository
}

func NewReactionUsecase(reactionRepo domain.ReactionRepository, userRepo domain.UserRepository, postRepo domain.PostRepository, commentRepo domain.CommentRepository) ReactionUseCase {
	return &reactionUsecase{
		reactionRepo: reactionRepo,
		userRepo:     userRepo,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
	}
}

func (u *reactionUsecase) AddReaction(ctx context.Context, reaction *domain.Reaction) (*domain.Reaction, error) {

	_, err := pghelpers.FindEntity(ctx, u.userRepo, "id", reaction.UserID, "User not found")
	if err != nil {
		return nil, err
	}

	switch reaction.EntityType {
	case "post":
		_, err = pghelpers.FindEntity(ctx, u.postRepo, "id", reaction.EntityID, "Post not found")
		if err != nil {
			return nil, err
		}
	case "comment":
		_, err = pghelpers.FindEntity(ctx, u.commentRepo, "id", reaction.EntityID, "Comment not found")
		if err != nil {
			return nil, err
		}
	default:
		return nil, apperror.NewError(apperror.InvalidData, "Invalid entity type", nil, "")
	}

	reaction.ID = helpers.NewULID()

	err = u.reactionRepo.AddReaction(ctx, reaction)
	if err != nil {
		return nil, err
	}

	return reaction, nil
}
