package post

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/helpers"
	"strings"
	"time"
)

type SearchResult struct {
	Posts []*domain.ExtendedPost
	Total int
}

type PostUseCase interface {
	Create(ctx context.Context, post *domain.Post) (*domain.ExtendedPost, error)
	Get(ctx context.Context, id int) (*domain.ExtendedPost, error)
	Search(ctx context.Context, params dto.SearchPostsParams) (*SearchResult, error)
	GetInterestedPosts(ctx context.Context, params dto.InterestedPostsParams) (*SearchResult, error)
	AddComment(ctx context.Context, comment *domain.Comment) (*domain.CommentWithUser, error)
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

func buildCriteria(field string, values []int) *criteria.Criteria {
	return &criteria.Criteria{
		Filters: []criteria.Filter{
			{Field: field, Value: values, Operator: criteria.OperatorIn},
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
	posts []*domain.Post) ([]*domain.ExtendedPost, error) {
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

		user, err := helpers.FindEntity(ctx, p.userRepository, "id", post.CreatedBy, "User not found")
		if err != nil {
			return nil, err
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

	existingSpace.UpdatedAt = time.Now()
	existingSpace.UpdatedBy = post.CreatedBy
	existingSpace.Posts += 1
	if err := p.spaceRepository.Update(ctx, existingSpace); err != nil {
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
	extendedPosts, err := p.buildExtendedPosts(ctx, []*domain.Post{post})
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

	post, err := helpers.FindEntity(ctx, p.postRepository, "id", comment.PostID, "Post not found")
	if err != nil {
		return nil, err
	}
	post.UpdatedAt = time.Now()
	post.UpdatedBy = comment.CreatedBy
	if err := p.postRepository.Update(ctx, post); err != nil {
		return nil, err
	}

	space, err := helpers.FindEntity(ctx, p.spaceRepository, "id", post.SpaceID, "Space not found")
	if err != nil {
		return nil, err
	}
	space.UpdatedAt = time.Now()
	space.UpdatedBy = comment.CreatedBy
	if err := p.spaceRepository.Update(ctx, space); err != nil {
		return nil, err
	}

	return &domain.CommentWithUser{Comment: comment, User: user}, nil
}

func (p *postUseCase) Search(ctx context.Context, params dto.SearchPostsParams) (*SearchResult, error) {
	var searchQuery string
	if len(params.Query) > 2 {
		searchQuery = "%" + strings.TrimSpace(params.Query) + "%"
	} else {
		searchQuery = ""
	}

	var spaceID int
	if params.SpaceID > 0 {
		spaceID = params.SpaceID
	}

	sortDirection := criteria.OrderDirectionDesc
	if params.SortDirection == "asc" {
		sortDirection = criteria.OrderDirectionAsc
	}

	var logicalOp criteria.LogicalOperator
	if len(params.Query) > 0 {
		logicalOp = criteria.LogicalOperatorOr
	} else {
		logicalOp = criteria.LogicalOperatorAnd
	}

	searchCriteria := criteria.NewCriteriaBuilder().
		WithFilterAndCondition("title", searchQuery, criteria.OperatorILike, len(params.Query) > 0).
		WithFilterAndCondition("content", searchQuery, criteria.OperatorILike, len(params.Query) > 0).
		WithFilterAndCondition("space_id", spaceID, criteria.OperatorEqual, spaceID > 0).
		WithLogicalOperator(logicalOp).
		WithPagination(params.Page, params.PageSize).
		WithSort(params.OrderBy, sortDirection).
		Build()

	// Get total count without pagination
	countCriteria := criteria.NewCriteriaBuilder().
		WithFilterAndCondition("title", searchQuery, criteria.OperatorILike, len(params.Query) > 0).
		WithFilterAndCondition("content", searchQuery, criteria.OperatorILike, len(params.Query) > 0).
		WithFilterAndCondition("space_id", spaceID, criteria.OperatorEqual, spaceID > 0).
		WithLogicalOperator(logicalOp).
		Build()

	total, err := p.postRepository.Count(ctx, countCriteria)
	if err != nil {
		return nil, err
	}

	posts, err := p.postRepository.Search(ctx, searchCriteria)
	if err != nil {
		return nil, err
	}

	extendedPosts, err := p.buildExtendedPosts(ctx, posts)
	if err != nil {
		return nil, err
	}

	return &SearchResult{
		Posts: extendedPosts,
		Total: total,
	}, nil
}

func (p *postUseCase) GetInterestedPosts(ctx context.Context, params dto.InterestedPostsParams) (*SearchResult, error) {
	sortDirection := criteria.OrderDirectionDesc
	if params.SortDirection == "asc" {
		sortDirection = criteria.OrderDirectionAsc
	}

	spaceIDs, err := p.userSpaceRepository.FindSpacesIDsByUserID(ctx, params.UserID)
	if err != nil {
		return nil, err
	}

	if len(spaceIDs) == 0 {
		return &SearchResult{
			Posts: []*domain.ExtendedPost{},
			Total: 0,
		}, nil
	}

	searchCriteria := criteria.NewCriteriaBuilder().
		WithFilter("space_id", spaceIDs, criteria.OperatorIn).
		WithPagination(params.Page, params.PageSize).
		WithSort(params.OrderBy, sortDirection).
		Build()

	// Get total count without pagination
	countCriteria := criteria.NewCriteriaBuilder().
		WithFilter("space_id", spaceIDs, criteria.OperatorIn).
		Build()

	total, err := p.postRepository.Count(ctx, countCriteria)
	if err != nil {
		return nil, err
	}

	posts, err := p.postRepository.Search(ctx, searchCriteria)
	if err != nil {
		return nil, err
	}

	extendedPosts, err := p.buildExtendedPosts(ctx, posts)
	if err != nil {
		return nil, err
	}

	return &SearchResult{
		Posts: extendedPosts,
		Total: total,
	}, nil
}
