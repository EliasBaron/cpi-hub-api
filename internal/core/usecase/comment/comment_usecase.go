package comment

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/pkg/apperror"
	"time"
)

type CommentUseCase interface {
	Create(ctx context.Context, comment *domain.Comment) (*domain.CommentWithUser, error)
	Get(ctx context.Context, id int) (*domain.CommentWithUser, error)
}

type commentUseCase struct {
	commentRepository domain.CommentRepository
	userRepository    domain.UserRepository
	postRepository    domain.PostRepository
}

func NewCommentUsecase(commentRepo domain.CommentRepository, userRepo domain.UserRepository, postRepo domain.PostRepository) CommentUseCase {
	return &commentUseCase{
		commentRepository: commentRepo,
		userRepository:    userRepo,
		postRepository:    postRepo,
	}
}

func (c commentUseCase) Create(ctx context.Context, comment *domain.Comment) (*domain.CommentWithUser, error) {
	existingUser, err := c.userRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    comment.CreatedBy,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, apperror.NewNotFound("User not found", nil, "comment_usecase.go:Create")
	}

	existingPost, err := c.postRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    comment.PostID,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}
	if existingPost == nil {
		return nil, apperror.NewNotFound("Post not found", nil, "comment_usecase.go:Create")
	}

	comment.CreatedAt, comment.UpdatedAt = time.Now(), time.Now()
	comment.UpdatedBy = comment.CreatedBy

	err = c.commentRepository.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &domain.CommentWithUser{
		Comment: comment,
		User:    existingUser,
	}, nil
}

func (c commentUseCase) Get(ctx context.Context, id int) (*domain.CommentWithUser, error) {
	//TODO implement me
	panic("implement me")
}
