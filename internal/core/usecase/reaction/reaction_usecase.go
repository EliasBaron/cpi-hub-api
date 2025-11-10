package reaction

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/core/dto"
	pghelpers "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/helpers"
	"cpi-hub-api/pkg/apperror"
)

type ReactionUseCase interface {
	AddReaction(ctx context.Context, reaction *domain.Reaction) (*domain.Reaction, error)
	RemoveReaction(ctx context.Context, reactionID string) error
	GetLikesCount(ctx context.Context, getLikesCountDTO dto.GetLikesCountDTO) (*dto.LikesCountDTO, error)
	GetUserLikes(ctx context.Context, userID int, entitiesData dto.EntitiesDataDTO) ([]dto.UserLikeDTO, error)
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

	if !domain.IsValidEntityType(string(reaction.EntityType)) {
		return nil, apperror.NewError(apperror.InvalidData, "Invalid entity type", nil, "")
	}
	if !domain.IsValidActionType(string(reaction.Action)) {
		return nil, apperror.NewError(apperror.InvalidData, "Invalid action type", nil, "")
	}

	_, err := pghelpers.FindEntity(ctx, u.userRepo, "id", reaction.UserID, "User not found")
	if err != nil {
		return nil, err
	}

	switch reaction.EntityType {
	case domain.EntityTypePost:
		_, err = pghelpers.FindEntity(ctx, u.postRepo, "id", reaction.EntityID, "Post not found")
		if err != nil {
			return nil, err
		}
	case domain.EntityTypeComment:
		_, err = pghelpers.FindEntity(ctx, u.commentRepo, "id", reaction.EntityID, "Comment not found")
		if err != nil {
			return nil, err
		}
	default:
		return nil, apperror.NewError(apperror.InvalidData, "Invalid entity type", nil, "")
	}

	criteria := &criteria.Criteria{
		Filters: []criteria.Filter{
			{Field: "user_id", Operator: criteria.OperatorEqual, Value: reaction.UserID},
			{Field: "entity_type", Operator: criteria.OperatorEqual, Value: string(reaction.EntityType)},
			{Field: "entity_id", Operator: criteria.OperatorEqual, Value: reaction.EntityID},
		},
	}
	existingReaction, err := u.reactionRepo.FindReaction(ctx, criteria)

	if err != nil {
		return nil, err
	}

	if existingReaction != nil {
		reaction.ID = existingReaction.ID
		err = u.reactionRepo.UpdateReaction(ctx, reaction)
		if err != nil {
			return nil, err
		}
	} else {
		err = u.reactionRepo.AddReaction(ctx, reaction)
		if err != nil {
			return nil, err
		}
	}

	return reaction, nil
}

func (u *reactionUsecase) RemoveReaction(ctx context.Context, reactionID string) error {
	err := u.reactionRepo.DeleteReaction(ctx, reactionID)
	if err != nil {
		return err
	}
	return nil
}

func (u *reactionUsecase) GetLikesCount(ctx context.Context, getLikesCountDTO dto.GetLikesCountDTO) (*dto.LikesCountDTO, error) {

	buildBaseCriteria := func() *criteria.CriteriaBuilder {
		builder := criteria.NewCriteriaBuilder()
		builder.WithFilterAndCondition("entity_type", getLikesCountDTO.EntityType, criteria.OperatorEqual, getLikesCountDTO.EntityType != nil).
			WithFilterAndCondition("entity_id", getLikesCountDTO.EntityID, criteria.OperatorEqual, getLikesCountDTO.EntityID != nil).
			WithFilterAndCondition("user_id", getLikesCountDTO.UserID, criteria.OperatorEqual, getLikesCountDTO.UserID != nil)
		return builder
	}

	critLikes := buildBaseCriteria().WithFilter("action", string(domain.ActionTypeLike), criteria.OperatorEqual).Build()
	likesCount, err := u.reactionRepo.CountReactions(ctx, critLikes)
	if err != nil {
		return nil, err
	}

	critDislikes := buildBaseCriteria().WithFilter("action", string(domain.ActionTypeDislike), criteria.OperatorEqual).Build()
	dislikesCount, err := u.reactionRepo.CountReactions(ctx, critDislikes)
	if err != nil {
		return nil, err
	}

	var likesCountDTO dto.LikesCountDTO

	if getLikesCountDTO.UserID != nil {
		likesCountDTO.UserID = getLikesCountDTO.UserID
	}
	if getLikesCountDTO.EntityType != nil && getLikesCountDTO.EntityID != nil {
		likesCountDTO.EntityID = getLikesCountDTO.EntityID
		likesCountDTO.EntityType = getLikesCountDTO.EntityType
	}
	likesCountDTO.LikesCount = likesCount
	likesCountDTO.DislikesCount = dislikesCount

	return &likesCountDTO, nil
}

func (u *reactionUsecase) GetUserLikes(ctx context.Context, userID int, entitiesData dto.EntitiesDataDTO) ([]dto.UserLikeDTO, error) {
	result := make([]dto.UserLikeDTO, 0, len(entitiesData.Entities))

	for _, entity := range entitiesData.Entities {

		_, err := pghelpers.FindEntity(ctx, u.userRepo, "id", userID, "User not found")
		if err != nil {
			return nil, err
		}

		switch entity.EntityType {
		case domain.EntityTypePost:
			_, err = pghelpers.FindEntity(ctx, u.postRepo, "id", entity.EntityID, "Post not found")
			if err != nil {
				return nil, err
			}
		case domain.EntityTypeComment:
			_, err = pghelpers.FindEntity(ctx, u.commentRepo, "id", entity.EntityID, "Comment not found")
			if err != nil {
				return nil, err
			}
		default:
			return nil, apperror.NewError(apperror.InvalidData, "Invalid entity type", nil, "")
		}

		criteria := criteria.NewCriteriaBuilder().
			WithFilter("user_id", userID, criteria.OperatorEqual).
			WithFilter("entity_type", entity.EntityType, criteria.OperatorEqual).
			WithFilter("entity_id", entity.EntityID, criteria.OperatorEqual).
			Build()

		reaction, err := u.reactionRepo.FindReaction(ctx, criteria)
		if err != nil {
			return nil, err
		}

		userLike := dto.UserLikeDTO{
			EntityType: string(entity.EntityType),
			EntityID:   entity.EntityID,
			Liked:      false,
			Disliked:   false,
		}

		if reaction != nil {
			switch reaction.Action {
			case domain.ActionTypeLike:
				userLike.Liked = true
			case domain.ActionTypeDislike:
				userLike.Disliked = true
			}
		}

		result = append(result, userLike)
	}

	return result, nil
}
