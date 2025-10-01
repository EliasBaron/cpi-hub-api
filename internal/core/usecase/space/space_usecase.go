package space

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/helpers"
	"cpi-hub-api/pkg/apperror"
	"strconv"
	"time"
)

type SpaceUseCase interface {
	Create(ctx context.Context, space *domain.Space) (*domain.SpaceWithUserAndCounts, error)
	Get(ctx context.Context, id string) (*domain.SpaceWithUserAndCounts, error)
	Search(ctx context.Context, criteria *domain.SpaceSearchCriteria) (*domain.SearchResult, error)
	GetUsersBySpace(ctx context.Context, spaceID string) ([]*domain.User, error)
}

type spaceUseCase struct {
	spaceRepository     domain.SpaceRepository
	userRepository      domain.UserRepository
	userSpaceRepository domain.UserSpaceRepository
	postRepository      domain.PostRepository
}

func NewSpaceUsecase(spaceRepository domain.SpaceRepository, userRepository domain.UserRepository, userSpaceRepository domain.UserSpaceRepository, postRepository domain.PostRepository) SpaceUseCase {
	return &spaceUseCase{
		spaceRepository:     spaceRepository,
		userRepository:      userRepository,
		userSpaceRepository: userSpaceRepository,
		postRepository:      postRepository,
	}
}

func (s *spaceUseCase) makeSpacesWithUsers(ctx context.Context, spaces []*domain.Space) ([]*domain.SpaceWithUserAndCounts, error) {
	var spacesWithUsers []*domain.SpaceWithUserAndCounts

	for _, space := range spaces {
		user, err := helpers.FindEntity(ctx, s.userRepository, "id", space.CreatedBy, "User not found")
		if err != nil {
			return nil, err
		}

		counts, err := s.getSpaceCounts(ctx, *space)
		if err != nil {
			return nil, err
		}

		spacesWithUsers = append(spacesWithUsers, &domain.SpaceWithUserAndCounts{
			Space:       space,
			User:        user,
			SpaceCounts: counts,
		})
	}

	return spacesWithUsers, nil
}

func (s *spaceUseCase) getSpaceCounts(ctx context.Context, space domain.Space) (domain.SpaceCounts, error) {
	var counts domain.SpaceCounts

	criteriaBuilder := criteria.NewCriteriaBuilder().
		WithFilter("space_id", space.ID, criteria.OperatorEqual)

	criteria := criteriaBuilder.Build()

	usersCount, err := s.userSpaceRepository.Count(ctx, criteria)
	if err != nil {
		return domain.SpaceCounts{}, err
	}

	postCount, err := s.postRepository.Count(ctx, criteria)
	if err != nil {
		return domain.SpaceCounts{}, err
	}

	counts.Users = usersCount
	counts.Posts = postCount

	return counts, nil
}

func (s *spaceUseCase) Create(ctx context.Context, space *domain.Space) (*domain.SpaceWithUserAndCounts, error) {
	existingUser, err := s.userRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    space.CreatedBy,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, apperror.NewNotFound("User not found", nil, "space_usecase.go:Create")
	}

	existingSpace, err := s.spaceRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "name",
				Value:    space.Name,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if existingSpace != nil {
		return nil, apperror.NewInvalidData("Space with this name already exists", nil, "space_usecase.go:Create")
	}

	space.CreatedAt, space.UpdatedAt = time.Now(), time.Now()
	space.UpdatedBy, space.CreatedBy = existingUser.ID, existingUser.ID

	err = s.spaceRepository.Create(ctx, space)
	if err != nil {
		return nil, err
	}

	err = s.userSpaceRepository.Update(ctx, existingUser.ID, []int{space.ID}, domain.AddUserToSpace)
	if err != nil {
		return nil, err
	}

	return &domain.SpaceWithUserAndCounts{
		Space: space,
		User:  existingUser,
		SpaceCounts: domain.SpaceCounts{
			Users: 1,
			Posts: 0,
		},
	}, nil
}

func (s *spaceUseCase) Get(ctx context.Context, id string) (*domain.SpaceWithUserAndCounts, error) {
	space, err := helpers.FindEntity(ctx, s.spaceRepository, "id", id, "Space not found")
	if err != nil {
		return nil, err
	}

	user, err := helpers.FindEntity(ctx, s.userRepository, "id", space.CreatedBy, "User not found")
	if err != nil {
		return nil, err
	}

	spaceCounts := domain.SpaceCounts{}
	spaceCounts, err = s.getSpaceCounts(ctx, *space)
	if err != nil {
		return nil, err
	}

	return &domain.SpaceWithUserAndCounts{
		Space: &domain.Space{
			ID:          space.ID,
			Name:        space.Name,
			Description: space.Description,
			CreatedAt:   space.CreatedAt,
			UpdatedAt:   space.UpdatedAt,
			CreatedBy:   space.CreatedBy,
			UpdatedBy:   space.UpdatedBy,
		},
		User: user,
		SpaceCounts: domain.SpaceCounts{
			Users: spaceCounts.Users,
			Posts: spaceCounts.Posts,
		},
	}, nil
}

func (s *spaceUseCase) Search(ctx context.Context, searchCriteria *domain.SpaceSearchCriteria) (*domain.SearchResult, error) {
	var direction criteria.Direction
	if searchCriteria.SortDirection == "asc" {
		direction = criteria.OrderDirectionAsc
	} else {
		direction = criteria.OrderDirectionDesc
	}

	criteriaBuilder := criteria.NewCriteriaBuilder().
		WithSort(searchCriteria.OrderBy, direction).
		WithPagination(searchCriteria.Page, searchCriteria.PageSize)

	if searchCriteria.Name != nil {
		criteriaBuilder.WithFilter("name", "%"+*searchCriteria.Name+"%", criteria.OperatorILike)
	}

	if searchCriteria.CreatedBy != nil {
		criteriaBuilder.WithFilter("created_by", *searchCriteria.CreatedBy, criteria.OperatorEqual)
	}

	criteria := criteriaBuilder.Build()

	totalCount, err := s.spaceRepository.Count(ctx, criteria)
	if err != nil {
		return nil, err
	}

	spaces, err := s.spaceRepository.FindAll(ctx, criteria)
	if err != nil {
		return nil, err
	}

	var spacesWithUsers []*domain.SpaceWithUserAndCounts
	if len(spaces) > 0 {
		spacesWithUsers, err = s.makeSpacesWithUsers(ctx, spaces)
		if err != nil {
			return nil, err
		}
	}

	return &domain.SearchResult{
		Data:     spacesWithUsers,
		Page:     searchCriteria.Page,
		PageSize: searchCriteria.PageSize,
		Total:    totalCount,
	}, nil
}

func (s *spaceUseCase) GetUsersBySpace(ctx context.Context, spaceID string) ([]*domain.User, error) {
	spaceIDInt, err := strconv.Atoi(spaceID)
	if err != nil {
		return nil, apperror.NewInvalidData("Invalid space ID format", err, "space_usecase.go:GetUsersBySpace")
	}

	_, err = helpers.FindEntity(ctx, s.spaceRepository, "id", spaceID, "Space not found")
	if err != nil {
		return nil, err
	}

	userIDs, err := s.userSpaceRepository.FindUserIDsBySpaceID(ctx, spaceIDInt)
	if err != nil {
		return nil, err
	}

	if len(userIDs) == 0 {
		return []*domain.User{}, nil
	}

	var users []*domain.User
	for _, userID := range userIDs {
		user, err := helpers.FindEntity(ctx, s.userRepository, "id", strconv.Itoa(userID), "User not found")
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
