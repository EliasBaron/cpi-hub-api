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
	SearchPosts(ctx context.Context, query string) ([]*domain.Post, error)
	GetPostsByUserSpaces(ctx context.Context, userId int) ([]*domain.ExtendedPost, error)
}

type postUseCase struct {
	postRepository      domain.PostRepository
	spaceRepository     domain.SpaceRepository
	userRepository      domain.UserRepository
	commentRepository   domain.CommentRepository
	userSpaceRepository domain.UserSpaceRepository
}

func NewPostUsecase(
	postRepo domain.PostRepository,
	spaceRepo domain.SpaceRepository,
	userRepo domain.UserRepository,
	commentRepo domain.CommentRepository,
	userSpaceRepo domain.UserSpaceRepository,
) PostUseCase {
	return &postUseCase{
		postRepository:      postRepo,
		spaceRepository:     spaceRepo,
		userRepository:      userRepo,
		commentRepository:   commentRepo,
		userSpaceRepository: userSpaceRepo,
	}
}

func intsToInterfaces(ints []int) []interface{} {
	res := make([]interface{}, len(ints))
	for i, v := range ints {
		res[i] = v
	}
	return res
}

func buildCriteria(field string, values []int) *criteria.Criteria {
	if len(values) == 1 {
		return &criteria.Criteria{
			Filters: []criteria.Filter{
				{Field: field, Value: values[0], Operator: criteria.OperatorEqual},
			},
		}
	}
	return &criteria.Criteria{
		Filters: []criteria.Filter{
			{Field: field, Value: intsToInterfaces(values), Operator: criteria.OperatorIn},
		},
	}
}

func (p *postUseCase) getCommentsWithUsers(ctx context.Context, postIDs []int) (map[int][]*domain.CommentWithUser, error) {
	comments, err := p.commentRepository.FindAll(ctx, buildCriteria("post_id", postIDs))
	if err != nil {
		return nil, err
	}

	commentsMap := make(map[int][]*domain.CommentWithUser)
	userCache := make(map[int]*domain.User)

	for _, comment := range comments {
		user := userCache[comment.CreatedBy]
		if user == nil {
			user, err = helpers.FindEntity(ctx, p.userRepository, "id", comment.CreatedBy, "User not found for comment")
			if err != nil {
				return nil, err
			}
			userCache[comment.CreatedBy] = user
		}
		commentsMap[comment.PostID] = append(commentsMap[comment.PostID], &domain.CommentWithUser{
			Comment: comment,
			User:    user,
		})
	}
	return commentsMap, nil
}

func (p *postUseCase) buildExtendedPosts(
	ctx context.Context,
	posts []*domain.Post,
	user *domain.User,
) ([]*domain.ExtendedPost, error) {
	if len(posts) == 0 {
		return []*domain.ExtendedPost{}, nil
	}

	postIDs := make([]int, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}
	commentsMap, err := p.getCommentsWithUsers(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	spaceCache := make(map[int]*domain.Space)
	result := make([]*domain.ExtendedPost, 0, len(posts))

	for _, post := range posts {
		space := spaceCache[post.SpaceID]
		if space == nil {
			space, err = helpers.FindEntity(ctx, p.spaceRepository, "id", post.SpaceID, "Space not found")
			if err != nil {
				return nil, err
			}
			spaceCache[post.SpaceID] = space
		}

		result = append(result, &domain.ExtendedPost{
			Post:     post,
			Space:    space,
			User:     user,
			Comments: commentsMap[post.ID],
		})
	}
	return result, nil
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
	user, err := helpers.FindEntity(ctx, p.userRepository, "id", post.CreatedBy, "User not found")
	if err != nil {
		return nil, err
	}
	extendedPosts, err := p.buildExtendedPosts(ctx, []*domain.Post{post}, user)
	if err != nil {
		return nil, err
	}
	return extendedPosts[0], nil
}

func (p *postUseCase) AddComment(ctx context.Context, comment *domain.Comment) (*domain.CommentWithUser, error) {
	user, err := helpers.FindEntity(ctx, p.userRepository, "id", comment.CreatedBy, "User not found")
	if err != nil {
		return nil, err
	}

	comment.CreatedAt, comment.UpdatedAt = time.Now(), time.Now()
	comment.UpdatedBy = comment.CreatedBy

	if err := p.commentRepository.Create(ctx, comment); err != nil {
		return nil, err
	}
	return &domain.CommentWithUser{Comment: comment, User: user}, nil
}

func (p *postUseCase) SearchPosts(ctx context.Context, query string) ([]*domain.Post, error) {
	return p.postRepository.SearchByTitleOrContent(ctx, query)
}

func (p *postUseCase) GetPostsByUserSpaces(ctx context.Context, userId int) ([]*domain.ExtendedPost, error) {
	user, err := helpers.FindEntity(ctx, p.userRepository, "id", userId, "User not found")
	if err != nil {
		return nil, err
	}

	userSpacesIds, err := p.userSpaceRepository.FindSpacesIDsByUser(ctx, userId)
	if err != nil || len(userSpacesIds) == 0 {
		return []*domain.ExtendedPost{}, err
	}

	posts, err := p.postRepository.FindAll(ctx, buildCriteria("space_id", userSpacesIds))
	if err != nil {
		return nil, err
	}

	return p.buildExtendedPosts(ctx, posts, user)
}
