package post

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/pkg/apperror"
	"time"
)

type PostUseCase interface {
	Create(ctx context.Context, post *domain.Post) (*domain.ExtendedPost, error)
	Get(ctx context.Context, id int) (*domain.ExtendedPost, error)
}

type postUseCase struct {
	postRepository    domain.PostRepository
	spaceRepository   domain.SpaceRepository
	userRepository    domain.UserRepository
	commentRepository domain.CommentRepository
}

func NewPostUsecase(postRepo domain.PostRepository, spaceRepo domain.SpaceRepository, userRepo domain.UserRepository, commentRepo domain.CommentRepository) PostUseCase {
	return &postUseCase{
		postRepository:    postRepo,
		spaceRepository:   spaceRepo,
		userRepository:    userRepo,
		commentRepository: commentRepo,
	}
}

func (p *postUseCase) Create(ctx context.Context, post *domain.Post) (*domain.ExtendedPost, error) {
	existingUser, err := p.userRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    post.CreatedBy,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, apperror.NewNotFound("User not found", nil, "post_usecase.go:Create")
	}

	existingSpace, err := p.spaceRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    post.SpaceID,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if existingSpace == nil {
		return nil, apperror.NewNotFound("Space not found", nil, "post_usecase.go:Create")
	}

	post.CreatedAt, post.UpdatedAt = time.Now(), time.Now()
	post.UpdatedBy = post.CreatedBy

	err = p.postRepository.Create(ctx, post)
	if err != nil {
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
	post, err := p.postRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    id,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if post == nil {
		return nil, apperror.NewNotFound("Post not found", nil, "post_usecase.go:Get")
	}

	space, err := p.spaceRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    post.SpaceID,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if space == nil {
		return nil, apperror.NewNotFound("Space not found", nil, "post_usecase.go:Get")
	}

	user, err := p.userRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    post.CreatedBy,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperror.NewNotFound("User not found", nil, "post_usecase.go:Get")
	}

	comments, err := p.commentRepository.FindAllByPostID(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	return &domain.ExtendedPost{
		Post:     post,
		Space:    space,
		User:     user,
		Comments: comments,
	}, nil
}
