package comment

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
)

type SearchResult struct {
	Comments []*domain.CommentWithInfo
	Total    int
}

type CommentUseCase interface {
	Search(ctx context.Context, params dto.SearchCommentsParams) (*SearchResult, error)
	Update(ctx context.Context, params dto.UpdateCommentDTO) error
}

type commentUseCase struct {
	commentRepository domain.CommentRepository
}

func NewCommentUsecase(commentRepo domain.CommentRepository) CommentUseCase {
	return &commentUseCase{
		commentRepository: commentRepo,
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
