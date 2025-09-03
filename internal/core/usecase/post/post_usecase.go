package post

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/helpers"
	"time"
)

type PostUseCase interface {
	Create(ctx context.Context, post *domain.Post) (*domain.ExtendedPost, error)
	Get(ctx context.Context, id int) (*domain.ExtendedPost, error)
	AddComment(ctx context.Context, comment *domain.Comment) (*domain.CommentWithUser, error)
}

type postUseCase struct {
	postRepository    domain.PostRepository
	spaceRepository   domain.SpaceRepository
	userRepository    domain.UserRepository
	commentRepository domain.CommentRepository
}

func NewPostUsecase(
	postRepo domain.PostRepository,
	spaceRepo domain.SpaceRepository,
	userRepo domain.UserRepository,
	commentRepo domain.CommentRepository,
) PostUseCase {
	return &postUseCase{
		postRepository:    postRepo,
		spaceRepository:   spaceRepo,
		userRepository:    userRepo,
		commentRepository: commentRepo,
	}
}

func (p *postUseCase) Create(ctx context.Context, post *domain.Post) (*domain.ExtendedPost, error) {
	existingUser, err := helpers.FindEntity(ctx, p.userRepository, "id", post.CreatedBy, "User not found")
	if err != nil {
		return nil, err
	}

	existingSpace, err := helpers.FindEntity(ctx, p.spaceRepository, "id", post.SpaceID, "Space not found")
	if err != nil {
		return nil, err
	}

	post.CreatedAt, post.UpdatedAt = time.Now(), time.Now()
	post.UpdatedBy = post.CreatedBy

	if err := p.postRepository.Create(ctx, post); err != nil {
		return nil, err
	}

	return &domain.ExtendedPost{
		Post:     post,
		Space:    existingSpace,
		User:     existingUser,
		Comments: []*domain.CommentWithUser{},
	}, nil
}

func (p *postUseCase) Get(ctx context.Context, id int) (*domain.ExtendedPost, error) {
	post, err := helpers.FindEntity(ctx, p.postRepository, "id", id, "Post not found")
	if err != nil {
		return nil, err
	}

	space, err := helpers.FindEntity(ctx, p.spaceRepository, "id", post.SpaceID, "Space not found")
	if err != nil {
		return nil, err
	}

	user, err := helpers.FindEntity(ctx, p.userRepository, "id", post.CreatedBy, "User not found")
	if err != nil {
		return nil, err
	}

	comments, err := p.commentRepository.FindAll(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{Field: "post_id", Value: post.ID, Operator: criteria.OperatorEqual},
		},
	})
	if err != nil {
		return nil, err
	}

	commentsWithUsers := make([]*domain.CommentWithUser, 0, len(comments))
	for _, comment := range comments {
		commentUser, err := helpers.FindEntity(ctx, p.userRepository, "id", comment.CreatedBy, "User not found for comment")
		if err != nil {
			return nil, err
		}
		commentsWithUsers = append(commentsWithUsers, &domain.CommentWithUser{
			Comment: comment,
			User:    commentUser,
		})
	}

	return &domain.ExtendedPost{
		Post:     post,
		Space:    space,
		User:     user,
		Comments: commentsWithUsers,
	}, nil
}

func (c postUseCase) AddComment(ctx context.Context, comment *domain.Comment) (*domain.CommentWithUser, error) {

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
