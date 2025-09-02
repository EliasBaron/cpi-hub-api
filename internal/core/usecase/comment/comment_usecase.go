package comment

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/pkg/helpers"
	"time"
)

type CommentUseCase interface {
	Create(ctx context.Context, comment *domain.Comment) (*domain.CommentWithUser, error)
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

	user, err := helpers.FindEntity(ctx, c.userRepository, "id", comment.CreatedBy, "User not found")
	if err != nil {
		return nil, err
	}

	comment.CreatedAt, comment.UpdatedAt = time.Now(), time.Now()
	comment.UpdatedBy = comment.CreatedBy

	err = c.commentRepository.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &domain.CommentWithUser{
		Comment: comment,
		User:    user,
	}, nil
}
