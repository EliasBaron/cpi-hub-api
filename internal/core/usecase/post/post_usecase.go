package post

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/events"
	pghelpers "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/helpers"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
	"fmt"
	"log"
	"strings"
)

type SearchResult struct {
	Posts []*domain.ExtendedPost
	Total int
}

type PostUseCase interface {
	Create(ctx context.Context, post *domain.Post) (*domain.ExtendedPost, error)
	Get(ctx context.Context, id int) (*domain.ExtendedPost, error)
	GetTrendingPosts(ctx context.Context, params dto.TrendingPostsParams) (*SearchResult, error)
	Search(ctx context.Context, params dto.SearchPostsParams) (*SearchResult, error)
	GetInterestedPosts(ctx context.Context, params dto.InterestedPostsParams) (*SearchResult, error)
	AddComment(ctx context.Context, commentDTO dto.CreateComment) (*domain.CommentWithInfo, error)
	Update(ctx context.Context, updatePostDTO *dto.UpdatePost) error
	Delete(ctx context.Context, postID int) error
}

type postUseCase struct {
	postRepository      domain.PostRepository
	spaceRepository     domain.SpaceRepository
	userRepository      domain.UserRepository
	commentRepository   domain.CommentRepository
	userSpaceRepository domain.UserSpaceRepository
	reactionRepository  domain.ReactionRepository
	eventEmitter        events.EventEmitter
}

func NewPostUsecase(
	postRepo domain.PostRepository,
	spaceRepo domain.SpaceRepository,
	userRepo domain.UserRepository,
	commentRepo domain.CommentRepository,
	userSpaceRepo domain.UserSpaceRepository,
	reactionRepo domain.ReactionRepository,
	eventEmitter events.EventEmitter,
) PostUseCase {
	return &postUseCase{
		postRepository:      postRepo,
		spaceRepository:     spaceRepo,
		userRepository:      userRepo,
		commentRepository:   commentRepo,
		userSpaceRepository: userSpaceRepo,
		reactionRepository:  reactionRepo,
		eventEmitter:        eventEmitter,
	}
}

func (p *postUseCase) getCommentsWithUsers(ctx context.Context, postIDs []int) (map[int][]*domain.CommentWithInfo, error) {
	comments, err := p.commentRepository.FindAll(ctx, criteria.NewCriteriaBuilder().
		WithFilter("post_id", postIDs, criteria.OperatorIn).
		WithSort("created_at", criteria.OrderDirectionDesc).
		Build())

	if err != nil {
		return nil, err
	}

	commentsMap := make(map[int][]*domain.CommentWithInfo)
	for _, comment := range comments {
		commentsMap[comment.Comment.PostID] = append(commentsMap[comment.Comment.PostID], comment)
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
	userCache := make(map[int]*domain.User)
	result := make([]*domain.ExtendedPost, 0, len(posts))

	for _, post := range posts {
		space := spaceCache[post.SpaceID]
		if space == nil {
			space, err = pghelpers.FindEntity(ctx, p.spaceRepository, "id", post.SpaceID, "Space not found")
			if err != nil {
				return nil, err
			}
			spaceCache[post.SpaceID] = space
		}

		user := userCache[post.CreatedBy]
		if user == nil {
			user, err = pghelpers.FindEntity(ctx, p.userRepository, "id", post.CreatedBy, "User not found")
			if err != nil {
				return nil, err
			}
			userCache[post.CreatedBy] = user
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
	existingUser, err := pghelpers.FindEntity(ctx, p.userRepository, "id", post.CreatedBy, "User not found")
	if err != nil {
		return nil, err
	}
	existingSpace, err := pghelpers.FindEntity(ctx, p.spaceRepository, "id", post.SpaceID, "Space not found")
	if err != nil {
		return nil, err
	}

	post.CreatedAt, post.UpdatedAt = helpers.GetTime(), helpers.GetTime()
	post.UpdatedBy = post.CreatedBy

	if err := p.postRepository.Create(ctx, post); err != nil {
		return nil, err
	}

	existingSpace.UpdatedAt = helpers.GetTime()
	existingSpace.UpdatedBy = post.CreatedBy
	if err := p.spaceRepository.Update(ctx, existingSpace); err != nil {
		return nil, err
	}

	return &domain.ExtendedPost{
		Post:     post,
		Space:    existingSpace,
		User:     existingUser,
		Comments: []*domain.CommentWithInfo{},
	}, nil
}

func (p *postUseCase) Get(ctx context.Context, id int) (*domain.ExtendedPost, error) {
	post, err := pghelpers.FindEntity(ctx, p.postRepository, "id", id, "Post not found")
	if err != nil {
		return nil, err
	}
	extendedPosts, err := p.buildExtendedPosts(ctx, []*domain.Post{post})
	if err != nil {
		return nil, err
	}
	return extendedPosts[0], nil
}

func (p *postUseCase) AddComment(ctx context.Context, commentDTO dto.CreateComment) (*domain.CommentWithInfo, error) {
	comment := commentDTO.ToDomain()

	user, err := pghelpers.FindEntity(ctx, p.userRepository, "id", comment.CreatedBy, "User not found")
	if err != nil {
		return nil, err
	}

	comment.CreatedAt, comment.UpdatedAt = helpers.GetTime(), helpers.GetTime()

	if comment.ParentID != nil && *comment.ParentID > 0 {
		parentComment, err := pghelpers.FindEntity(ctx, p.commentRepository, "id", *comment.ParentID, "Parent comment not found")
		if err != nil {
			return nil, err
		}
		if parentComment.Comment.ParentID != nil {
			return nil, apperror.NewInvalidData("You can only reply to root comments (comments without a parent)", nil, "post_usecase.go:AddComment")
		}
	}

	if err := p.commentRepository.Create(ctx, comment); err != nil {
		return nil, err
	}

	post, err := pghelpers.FindEntity(ctx, p.postRepository, "id", comment.PostID, "Post not found")
	if err != nil {
		return nil, err
	}
	post.UpdatedAt = helpers.GetTime()
	post.UpdatedBy = comment.CreatedBy
	if err := p.postRepository.Update(ctx, post); err != nil {
		return nil, err
	}

	space, err := pghelpers.FindEntity(ctx, p.spaceRepository, "id", post.SpaceID, "Space not found")
	if err != nil {
		return nil, err
	}
	space.UpdatedAt = helpers.GetTime()
	space.UpdatedBy = comment.CreatedBy
	if err := p.spaceRepository.Update(ctx, space); err != nil {
		return nil, err
	}

	if p.eventEmitter != nil {
		var event *domain.Event

		if comment.ParentID != nil && *comment.ParentID > 0 {
			parentComment, err := pghelpers.FindEntity(ctx, p.commentRepository, "id", *comment.ParentID, "Parent comment not found")
			if err == nil && parentComment != nil {
				if parentComment.Comment.CreatedBy != comment.CreatedBy {
					event = &domain.Event{
						Type:         "comment_reply_created",
						UserID:       comment.CreatedBy,
						TargetUserID: parentComment.Comment.CreatedBy,
						Metadata: map[string]interface{}{
							"post_id":           comment.PostID,
							"comment_id":        comment.ID,
							"parent_comment_id": *comment.ParentID,
							"comment_content":   comment.Content,
							"is_reply":          true,
						},
						Timestamp: comment.CreatedAt,
					}
				}
			}
		} else {
			if post.CreatedBy != comment.CreatedBy {
				event = &domain.Event{
					Type:         "comment_created",
					UserID:       comment.CreatedBy,
					TargetUserID: post.CreatedBy,
					Metadata: map[string]interface{}{
						"post_id":         comment.PostID,
						"comment_id":      comment.ID,
						"comment_content": comment.Content,
						"is_reply":        false,
					},
					Timestamp: comment.CreatedAt,
				}
			}
		}

		if event != nil {
			err := p.eventEmitter.EmitEvent(ctx, event)
			if err != nil {
				log.Printf("Error emitting comment event: %v", err)
			}
		}
	}

	return &domain.CommentWithInfo{
		Comment: comment,
		User:    user,
		Space:   space,
	}, nil
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

	var userID int
	if params.UserID > 0 {
		userID = params.UserID
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
		WithFilterAndCondition("created_by", userID, criteria.OperatorEqual, userID > 0).
		WithLogicalOperator(logicalOp).
		WithPagination(params.Page, params.PageSize).
		WithSort(params.OrderBy, sortDirection).
		Build()

	countCriteria := criteria.NewCriteriaBuilder().
		WithFilterAndCondition("title", searchQuery, criteria.OperatorILike, len(params.Query) > 0).
		WithFilterAndCondition("content", searchQuery, criteria.OperatorILike, len(params.Query) > 0).
		WithFilterAndCondition("space_id", spaceID, criteria.OperatorEqual, spaceID > 0).
		WithFilterAndCondition("created_by", userID, criteria.OperatorEqual, userID > 0).
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

func (p *postUseCase) GetTrendingPosts(ctx context.Context, params dto.TrendingPostsParams) (*SearchResult, error) {
	since, err := helpers.ParseTimeFrame(params.TimeFrame)
	if err != nil {
		return nil, apperror.NewInvalidData(fmt.Sprintf("Invalid time_frame: %s", params.TimeFrame), err, "post_usecase.go:GetTrendingPosts")
	}

	// Build criteria for aggregating reactions on posts since the timeframe
	reactionCriteria := criteria.NewCriteriaBuilder().
		WithFilter("entity_type", string(domain.EntityTypePost), criteria.OperatorEqual).
		WithFilter("action", string(domain.ActionTypeLike), criteria.OperatorEqual).
		WithFilter("timestamp", since, criteria.OperatorGte).
		WithPagination(params.Page, params.PageSize).
		Build()

	// Get top posts by reaction count, grouped by entity_id (post_id)
	topReactions, total, err := p.reactionRepository.GetTopReactionEntities(ctx, reactionCriteria, "entity_id")
	if err != nil {
		return nil, err
	}

	if len(topReactions) == 0 {
		return &SearchResult{Posts: []*domain.ExtendedPost{}, Total: 0}, nil
	}

	// Extract post IDs from aggregation results, preserving order
	postIDs := make([]int, len(topReactions))
	for i, tr := range topReactions {
		postIDs[i] = tr.EntityID
	}

	// Fetch posts from Postgres by IDs
	postsCriteria := criteria.NewCriteriaBuilder().
		WithFilter("id", postIDs, criteria.OperatorIn).
		Build()

	posts, err := p.postRepository.Search(ctx, postsCriteria)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup and preserve trending order
	postMap := make(map[int]*domain.Post)
	for _, post := range posts {
		postMap[post.ID] = post
	}

	// Rebuild posts slice in the order of topReactions (most likes first)
	orderedPosts := make([]*domain.Post, 0, len(topReactions))
	for _, tr := range topReactions {
		if post, exists := postMap[tr.EntityID]; exists {
			orderedPosts = append(orderedPosts, post)
		}
	}

	// Build extended posts with user, space, and comments
	extendedPosts, err := p.buildExtendedPosts(ctx, orderedPosts)
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

	searchCriteria := criteria.NewCriteriaBuilder().
		WithFilter("space_id", spaceIDs, criteria.OperatorIn).
		WithPagination(params.Page, params.PageSize).
		WithSort(params.OrderBy, sortDirection).
		Build()

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

func (p *postUseCase) Update(ctx context.Context, updatePostDTO *dto.UpdatePost) error {
	existingPost, err := pghelpers.FindEntity(ctx, p.postRepository, "id", updatePostDTO.PostID, "Post not found")
	if err != nil {
		return err
	}

	if updatePostDTO.Title != "" {
		existingPost.Title = updatePostDTO.Title
	}
	if updatePostDTO.Content != "" {
		existingPost.Content = updatePostDTO.Content
	}
	existingPost.UpdatedAt = helpers.GetTime()

	if err := p.postRepository.Update(ctx, existingPost); err != nil {
		return err
	}

	return nil
}

func (p *postUseCase) Delete(ctx context.Context, postID int) error {
	existingPost, err := pghelpers.FindEntity(ctx, p.postRepository, "id", postID, "Post not found")
	if err != nil {
		return err
	}
	if existingPost == nil {
		return apperror.NewNotFound("post not found", nil, "post_usecase.go:Delete")
	}
	extendedPosts, err := p.buildExtendedPosts(ctx, []*domain.Post{existingPost})
	if err != nil {
		return err
	}
	existingPost = extendedPosts[0].Post

	if err := p.postRepository.Delete(ctx, existingPost.ID); err != nil {
		return err
	}

	return nil
}
