package comment

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
	"fmt"
)

type SearchResult struct {
	Comments []*domain.CommentWithInfo
	Total    int
}

type CommentUseCase interface {
	Search(ctx context.Context, params dto.SearchCommentsParams) (*SearchResult, error)
	GetTrendingComments(ctx context.Context, params dto.TrendingCommentsParams) (*SearchResult, error)
	Update(ctx context.Context, params dto.UpdateCommentDTO) error
	Delete(ctx context.Context, commentID int) error
}

type commentUseCase struct {
	commentRepository  domain.CommentRepository
	reactionRepository domain.ReactionRepository
}

func NewCommentUsecase(commentRepo domain.CommentRepository, reactionRepo domain.ReactionRepository) CommentUseCase {
	return &commentUseCase{
		commentRepository:  commentRepo,
		reactionRepository: reactionRepo,
	}
}

func (c *commentUseCase) Search(ctx context.Context, params dto.SearchCommentsParams) (*SearchResult, error) {
	sortDirection := criteria.OrderDirectionDesc
	if params.SortDirection == "asc" {
		sortDirection = criteria.OrderDirectionAsc
	}

	builder := criteria.NewCriteriaBuilder()

	if params.UserID != nil && *params.UserID > 0 {
		builder.WithFilter("created_by", *params.UserID, criteria.OperatorEqual)
	}

	if params.PostID != nil && *params.PostID > 0 {
		builder.WithFilter("post_id", *params.PostID, criteria.OperatorEqual)
	}

	total, err := c.commentRepository.Count(ctx, builder.Build())
	if err != nil {
		return nil, err
	}

	commentsWithSpace, err := c.commentRepository.FindAll(ctx, builder.
		WithPagination(params.Page, params.PageSize).
		WithSort(params.OrderBy, sortDirection).
		Build())

	if err != nil {
		return nil, err
	}

	return &SearchResult{
		Comments: commentsWithSpace,
		Total:    total,
	}, nil
}

func (c *commentUseCase) Update(ctx context.Context, params dto.UpdateCommentDTO) error {

	searchCriteria := criteria.NewCriteriaBuilder().
		WithFilter("id", params.CommentID, criteria.OperatorEqual).
		Build()

	existingComment, err := c.commentRepository.Find(ctx, searchCriteria)
	if err != nil {
		return err
	}

	if existingComment == nil {
		return apperror.NewNotFound("comment not found", nil, "comment_usecase.go:Update")
	}

	existingComment.Comment.Content = params.Content
	existingComment.Comment.UpdatedAt = helpers.GetTime()

	if err := c.commentRepository.Update(ctx, existingComment.Comment); err != nil {
		return err
	}

	return nil
}

func (c *commentUseCase) Delete(ctx context.Context, commentID int) error {
	searchCriteria := criteria.NewCriteriaBuilder().
		WithFilter("id", commentID, criteria.OperatorEqual).
		Build()

	existingComment, err := c.commentRepository.Find(ctx, searchCriteria)
	if err != nil {
		return err
	}
	if existingComment == nil {
		return apperror.NewNotFound("comment not found", nil, "comment_usecase.go:Delete")
	}

	if err := c.commentRepository.Delete(ctx, commentID); err != nil {
		return err
	}

	return nil
}

func (c *commentUseCase) GetTrendingComments(ctx context.Context, params dto.TrendingCommentsParams) (*SearchResult, error) {
	since, err := helpers.ParseTimeFrame(params.TimeFrame)
	if err != nil {
		return nil, apperror.NewInvalidData(fmt.Sprintf("Invalid time_frame: %s", params.TimeFrame), err, "comment_usecase.go:GetTrendingComments")
	}

	reactionCriteria := criteria.NewCriteriaBuilder().
		WithFilter("entity_type", string(domain.EntityTypeComment), criteria.OperatorEqual).
		WithFilter("action", string(domain.ActionTypeLike), criteria.OperatorEqual).
		WithFilter("timestamp", since, criteria.OperatorGte).
		WithPagination(params.Page, params.PageSize).
		Build()

	topReactions, total, err := c.reactionRepository.GetTopReactionEntities(ctx, reactionCriteria, "entity_id")
	if err != nil {
		return nil, err
	}

	if len(topReactions) == 0 {
		return &SearchResult{Comments: []*domain.CommentWithInfo{}, Total: 0}, nil
	}

	commentIDs := make([]int, len(topReactions))
	for i, tr := range topReactions {
		commentIDs[i] = tr.EntityID
	}

	commentsCriteria := criteria.NewCriteriaBuilder().
		WithFilter("id", commentIDs, criteria.OperatorIn).
		Build()

	comments, err := c.commentRepository.FindAll(ctx, commentsCriteria)
	if err != nil {
		return nil, err
	}

	commentMap := make(map[int]*domain.CommentWithInfo)
	for _, comment := range comments {
		commentMap[comment.Comment.ID] = comment
	}

	orderedComments := make([]*domain.CommentWithInfo, 0, len(topReactions))
	for _, tr := range topReactions {
		if comment, exists := commentMap[tr.EntityID]; exists {
			orderedComments = append(orderedComments, comment)
		}
	}

	return &SearchResult{
		Comments: orderedComments,
		Total:    total,
	}, nil
}
