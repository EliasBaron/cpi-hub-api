package user

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Get(ctx context.Context, id int) (*domain.UserWithSpaces, error)
	Update(ctx context.Context, dto dto.UpdateUserSpacesDTO) error
	GetSpacesByUser(ctx context.Context, userId int) ([]*domain.Space, error)
	Login(ctx context.Context, loginUser dto.LoginUser) (*domain.User, error)
	Search(ctx context.Context, params dto.SearchUsersParams) (*dto.PaginatedUsersResponse, error)
	UpdateUser(ctx context.Context, dto dto.UpdateUserDTO) (*domain.User, error)
}

type useCase struct {
	userRepository      domain.UserRepository
	spaceRepository     domain.SpaceRepository
	userSpaceRepository domain.UserSpaceRepository
}

func NewUserUsecase(userRepository domain.UserRepository, spaceRepository domain.SpaceRepository, userSpaceRepository domain.UserSpaceRepository) UserUseCase {
	return &useCase{
		userRepository:      userRepository,
		spaceRepository:     spaceRepository,
		userSpaceRepository: userSpaceRepository,
	}
}

func (u *useCase) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := u.userRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "email",
				Value:    user.Email,
				Operator: criteria.OperatorEqual,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, apperror.NewInvalidData("User with this email already exists", nil, "user_usecase.go:Create")
	}

	user.CreatedAt = helpers.GetTime()
	user.UpdatedAt = helpers.GetTime()

	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.NewInvalidData("Failed to hash password", err, "user_usecase.go:Create")
	}
	user.Password = string(cryptedPassword)

	err = u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *useCase) Get(ctx context.Context, id int) (*domain.UserWithSpaces, error) {
	user, err := u.userRepository.Find(ctx, &criteria.Criteria{
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

	if user == nil {
		return nil, apperror.NewNotFound("User not found", nil, "user_usecase.go:GetUserWithSpaces")
	}

	spaceIDs, err := u.userSpaceRepository.FindSpacesIDsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	spaces, err := u.spaceRepository.FindByIDs(ctx, spaceIDs)

	if err != nil {
		return nil, err
	}

	return &domain.UserWithSpaces{
		User:   user,
		Spaces: spaces,
	}, nil
}

func (u *useCase) GetSpacesByUser(ctx context.Context, userId int) ([]*domain.Space, error) {
	spaceIDs, err := u.userSpaceRepository.FindSpacesIDsByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return u.spaceRepository.FindByIDs(ctx, spaceIDs)
}

func (u *useCase) Update(ctx context.Context, dto dto.UpdateUserSpacesDTO) error {
	user, err := u.userRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    dto.UserID,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return err
	}

	if user == nil {
		return apperror.NewNotFound("User not found", nil, "user_usecase.go:Update")
	}

	if len(dto.SpaceIDs) == 0 {
		return apperror.NewInvalidData("Space IDs cannot be empty", nil, "user_usecase.go:Update")
	}

	for _, spaceID := range dto.SpaceIDs {
		exists, err := u.spaceRepository.Find(ctx, &criteria.Criteria{
			Filters: []criteria.Filter{
				{
					Field:    "id",
					Value:    spaceID,
					Operator: criteria.OperatorEqual,
				},
			},
		})
		if err != nil {
			return err
		}
		if exists == nil {
			return apperror.NewInvalidData("Space not found: "+strconv.Itoa(spaceID), nil, "user_usecase.go:Update")
		}
	}

	if err := u.userSpaceRepository.Update(ctx, user.ID, dto.SpaceIDs, dto.Action); err != nil {
		return err
	}

	return nil
}

func (u *useCase) Search(ctx context.Context, params dto.SearchUsersParams) (*dto.PaginatedUsersResponse, error) {
	builder := criteria.NewCriteriaBuilder()

	if params.FullName != "" {
		builder.WithFilter("CONCAT(name, ' ', last_name)", "%"+params.FullName+"%", criteria.OperatorILike)
	}

	builder.WithPagination(params.Page, params.PageSize)

	if params.OrderBy != "" {
		var direction criteria.Direction = criteria.OrderDirectionDesc
		if params.SortDirection == "asc" {
			direction = criteria.OrderDirectionAsc
		}
		builder.WithSort(params.OrderBy, direction)
	}

	searchCriteria := builder.Build()

	users, err := u.userRepository.Search(ctx, searchCriteria)
	if err != nil {
		return nil, err
	}

	countBuilder := criteria.NewCriteriaBuilder()
	if params.FullName != "" {
		countBuilder.WithFilter("CONCAT(name, ' ', last_name)", "%"+params.FullName+"%", criteria.OperatorILike)
	}
	countCriteria := countBuilder.Build()

	total, err := u.userRepository.Count(ctx, countCriteria)
	if err != nil {
		return nil, err
	}

	userDTOs := make([]dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = dto.ToUserDTO(user)
	}

	return &dto.PaginatedUsersResponse{
		Data:     userDTOs,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

func (u *useCase) Login(ctx context.Context, loginUser dto.LoginUser) (*domain.User, error) {
	user, err := u.userRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "email",
				Value:    loginUser.Email,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperror.NewInvalidData("Invalid email", nil, "user_usecase.go:Login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		return nil, apperror.NewInvalidData("Invalid password", nil, "user_usecase.go:Login")
	}

	return user, nil
}

func (u *useCase) UpdateUser(ctx context.Context, dto dto.UpdateUserDTO) (*domain.User, error) {
	user, err := u.userRepository.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{
				Field:    "id",
				Value:    dto.ID,
				Operator: criteria.OperatorEqual,
			},
		},
	})

	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NewNotFound("User not found", nil, "user_usecase.go:UpdateUser")
	}

	if dto.Name != nil {
		user.Name = *dto.Name
	}
	if dto.LastName != nil {
		user.LastName = *dto.LastName
	}
	if dto.Image != nil {
		user.Image = *dto.Image
	}
	user.UpdatedAt = helpers.GetTime()

	err = u.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
