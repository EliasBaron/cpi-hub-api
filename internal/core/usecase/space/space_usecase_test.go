package space

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/mock"
	"cpi-hub-api/pkg/apperror"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpaceRepository := mock.NewMockSpaceRepository(ctrl)
	mockUserRepository := mock.NewMockUserRepository(ctrl)
	mockUserSpaceRepository := mock.NewMockUserSpaceRepository(ctrl)
	mockPostRepository := mock.NewMockPostRepository(ctrl)

	spaceUseCase := NewSpaceUsecase(mockSpaceRepository, mockUserRepository, mockUserSpaceRepository, mockPostRepository)

	type args struct {
		context context.Context
		space   *domain.Space
	}

	type want struct {
		space *domain.SpaceWithUserAndCounts
		err   error
	}

	givenSpace := &domain.Space{
		ID:          1,
		Name:        "Test Space",
		Description: "Test Description",
		CreatedBy:   1,
	}

	givenUser := &domain.User{
		ID:   1,
		Name: "Test User",
	}

	tests := []struct {
		name  string
		args  args
		want  want
		calls []*gomock.Call
	}{
		{
			name: "success",
			args: args{
				context: context.Background(),
				space:   givenSpace,
			},
			want: want{
				space: &domain.SpaceWithUserAndCounts{
					Space: &domain.Space{
						ID:          1,
						Name:        "Test Space",
						Description: "Test Description",
					},
					User: &domain.User{
						ID:   1,
						Name: "Test User",
					},
					SpaceCounts: domain.SpaceCounts{
						Users: 1,
						Posts: 0,
					},
				},
			},
			calls: []*gomock.Call{
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenUser, nil),
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, nil),
				mockSpaceRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil),
				mockUserSpaceRepository.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
			},
		},
		{
			name: "error finding existing user",
			args: args{
				context: context.Background(),
				space:   givenSpace,
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error")),
			},
		},
		{
			name: "error user not found",
			args: args{
				context: context.Background(),
				space:   givenSpace,
			},
			want: want{
				err: apperror.NewNotFound("User not found", nil, "space_usecase.go:Create"),
			},
			calls: []*gomock.Call{
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, nil),
			},
		},
		{
			name: "error finding existing space",
			args: args{
				context: context.Background(),
				space:   givenSpace,
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenUser, nil),
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error")),
			},
		},
		{
			name: "error space already exists",
			args: args{
				context: context.Background(),
				space:   givenSpace,
			},
			want: want{
				err: apperror.NewInvalidData("Space with this name already exists", nil, "space_usecase.go:Create"),
			},
			calls: []*gomock.Call{
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenUser, nil),
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenSpace, nil),
			},
		},
		{
			name: "error creating space",
			args: args{
				context: context.Background(),
				space:   givenSpace,
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenUser, nil),
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, nil),
				mockSpaceRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("unexpected error")),
			},
		},
		{
			name: "error updating user space",
			args: args{
				context: context.Background(),
				space:   givenSpace,
			},
			want: want{
				err: errors.New("error updating user space"),
			},
			calls: []*gomock.Call{
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenUser, nil),
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, nil),
				mockSpaceRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil),
				mockUserSpaceRepository.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error updating user space")),
			},
		},
	}

	for _, test := range tests {
		calls := make([]interface{}, len(test.calls))
		for i, c := range test.calls {
			calls[i] = c
		}

		gomock.InOrder(calls...)

		got, gotErr := spaceUseCase.Create(test.args.context, test.args.space)

		assert.Equal(t, test.want.err, gotErr)
		if test.want.err == nil {
			assert.Equal(t, test.want.space.Space.ID, got.Space.ID)
			assert.Equal(t, test.want.space.User.ID, got.User.ID)
			assert.Equal(t, test.want.space.SpaceCounts.Users, got.SpaceCounts.Users)
			assert.Equal(t, test.want.space.SpaceCounts.Posts, got.SpaceCounts.Posts)
		} else {
			assert.Equal(t, test.want.err, gotErr)
		}
	}
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpaceRepository := mock.NewMockSpaceRepository(ctrl)
	mockUserRepository := mock.NewMockUserRepository(ctrl)
	mockUserSpaceRepository := mock.NewMockUserSpaceRepository(ctrl)
	mockPostRepository := mock.NewMockPostRepository(ctrl)

	spaceUseCase := NewSpaceUsecase(mockSpaceRepository, mockUserRepository, mockUserSpaceRepository, mockPostRepository)

	type args struct {
		context context.Context
		id      string
	}

	type want struct {
		space *domain.SpaceWithUserAndCounts
		err   error
	}

	givenSpace := &domain.Space{
		ID:          1,
		Name:        "Test Space",
		Description: "Test Description",
		CreatedBy:   1,
	}

	givenUser := &domain.User{
		ID:   1,
		Name: "Test User",
	}

	tests := []struct {
		name  string
		args  args
		want  want
		calls []*gomock.Call
	}{
		{
			name: "success",
			args: args{
				context: context.Background(),
				id:      "1",
			},
			want: want{
				space: &domain.SpaceWithUserAndCounts{
					Space: &domain.Space{
						ID:          1,
						Name:        "Test Space",
						Description: "Test Description",
					},
					User: &domain.User{
						ID:   1,
						Name: "Test User",
					},
					SpaceCounts: domain.SpaceCounts{
						Users: 1,
						Posts: 0,
					},
				},
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenSpace, nil),
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenUser, nil),
				mockUserSpaceRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(1, nil),
				mockPostRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(0, nil),
			},
		},
		{
			name: "error finding space",
			args: args{
				context: context.Background(),
				id:      "1",
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error")),
			},
		},
		{
			name: "error finding user",
			args: args{
				context: context.Background(),
				id:      "1",
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenSpace, nil),
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error")),
			},
		},
		{
			name: "error counting user spaces",
			args: args{
				context: context.Background(),
				id:      "1",
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenSpace, nil),
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenUser, nil),
				mockUserSpaceRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(0, errors.New("unexpected error")),
			},
		},
		{
			name: "error counting posts",
			args: args{
				context: context.Background(),
				id:      "1",
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenSpace, nil),
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenUser, nil),
				mockUserSpaceRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(0, nil),
				mockPostRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(0, errors.New("unexpected error")),
			},
		},
	}

	for _, test := range tests {
		calls := make([]interface{}, len(test.calls))
		for i, c := range test.calls {
			calls[i] = c
		}
		gomock.InOrder(calls...)

		got, gotErr := spaceUseCase.Get(test.args.context, test.args.id)

		assert.Equal(t, test.want.err, gotErr)
		if test.want.err == nil {
			assert.Equal(t, test.want.space.Space.ID, got.Space.ID)
			assert.Equal(t, test.want.space.User.ID, got.User.ID)
			assert.Equal(t, test.want.space.SpaceCounts.Users, got.SpaceCounts.Users)
			assert.Equal(t, test.want.space.SpaceCounts.Posts, got.SpaceCounts.Posts)
		} else {
			assert.Equal(t, test.want.err, gotErr)
		}
	}
}

func TestSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpaceRepository := mock.NewMockSpaceRepository(ctrl)
	mockUserRepository := mock.NewMockUserRepository(ctrl)
	mockUserSpaceRepository := mock.NewMockUserSpaceRepository(ctrl)
	mockPostRepository := mock.NewMockPostRepository(ctrl)

	spaceUseCase := NewSpaceUsecase(mockSpaceRepository, mockUserRepository, mockUserSpaceRepository, mockPostRepository)

	type args struct {
		context  context.Context
		criteria *domain.SpaceSearchCriteria
	}

	type want struct {
		result *domain.SearchResult
		err    error
	}

	givenSpace := &domain.Space{
		ID:          1,
		Name:        "Test Space",
		Description: "Test Description",
		CreatedBy:   1,
	}

	givenUser := &domain.User{
		ID:   1,
		Name: "Test User",
	}

	tests := []struct {
		name  string
		args  args
		want  want
		calls []*gomock.Call
	}{
		{
			name: "success",
			args: args{
				context: context.Background(),
				criteria: &domain.SpaceSearchCriteria{
					Query:         "Test Space",
					Name:          &givenSpace.Name,
					CreatedBy:     &givenSpace.CreatedBy,
					SortDirection: "asc",
				},
			},
			want: want{
				result: &domain.SearchResult{
					Data: []*domain.SpaceWithUserAndCounts{
						{Space: givenSpace, User: givenUser, SpaceCounts: domain.SpaceCounts{Users: 1, Posts: 0}},
					},
				},
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(1, nil),
				mockSpaceRepository.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return([]*domain.Space{givenSpace}, nil),
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenUser, nil),
				mockUserSpaceRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(1, nil),
				mockPostRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(0, nil),
			},
		},
		{
			name: "error counting spaces",
			args: args{
				context: context.Background(),
				criteria: &domain.SpaceSearchCriteria{
					Query: "Test Space",
				},
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(0, errors.New("unexpected error")),
			},
		},
		{
			name: "error counting user spaces",
			args: args{
				context: context.Background(),
				criteria: &domain.SpaceSearchCriteria{
					Query: "Test Space",
				},
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(1, nil),
				mockSpaceRepository.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return([]*domain.Space{givenSpace}, nil),
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenUser, nil),
				mockUserSpaceRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(0, errors.New("unexpected error")),
			},
		},
		{
			name: "error finding spaces",
			args: args{
				context: context.Background(),
				criteria: &domain.SpaceSearchCriteria{
					Name:      &givenSpace.Name,
					CreatedBy: &givenSpace.CreatedBy,
				},
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(1, nil),
				mockSpaceRepository.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error")),
			},
		},
		{
			name: "error finding user",
			args: args{
				context: context.Background(),
				criteria: &domain.SpaceSearchCriteria{
					Query: "Test Space",
				},
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Count(gomock.Any(), gomock.Any()).Return(1, nil),
				mockSpaceRepository.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return([]*domain.Space{givenSpace}, nil),
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error")),
			},
		},
	}

	for _, test := range tests {
		calls := make([]interface{}, len(test.calls))
		for i, c := range test.calls {
			calls[i] = c
		}
		gomock.InOrder(calls...)

		got, gotErr := spaceUseCase.Search(test.args.context, test.args.criteria)

		assert.Equal(t, test.want.err, gotErr)
		if test.want.err == nil {
			assert.Equal(t, test.want.result.Data, got.Data)
		} else {
			assert.Equal(t, test.want.err, gotErr)
		}
	}
}

func TestGetUsersBySpace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpaceRepository := mock.NewMockSpaceRepository(ctrl)
	mockUserRepository := mock.NewMockUserRepository(ctrl)
	mockUserSpaceRepository := mock.NewMockUserSpaceRepository(ctrl)
	mockPostRepository := mock.NewMockPostRepository(ctrl)

	spaceUseCase := NewSpaceUsecase(mockSpaceRepository, mockUserRepository, mockUserSpaceRepository, mockPostRepository)

	type args struct {
		context context.Context
		spaceID string
	}

	type want struct {
		users []*domain.User
		err   error
	}

	givenSpace := &domain.Space{
		ID:          1,
		Name:        "Test Space",
		Description: "Test Description",
		CreatedBy:   1,
	}

	givenUser := &domain.User{
		ID:   1,
		Name: "Test User",
	}

	tests := []struct {
		name  string
		args  args
		want  want
		calls []*gomock.Call
	}{
		{
			name: "success",
			args: args{
				context: context.Background(),
				spaceID: "1",
			},
			want: want{
				users: []*domain.User{givenUser},
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenSpace, nil),
				mockUserSpaceRepository.EXPECT().FindUserIDsBySpaceID(gomock.Any(), gomock.Any()).Return([]int{1}, nil),
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenUser, nil),
			},
		},
		{
			name: "success with no users",
			args: args{
				context: context.Background(),
				spaceID: "1",
			},
			want: want{
				users: []*domain.User{},
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenSpace, nil),
				mockUserSpaceRepository.EXPECT().FindUserIDsBySpaceID(gomock.Any(), gomock.Any()).Return([]int{}, nil),
			},
		},
		{
			name: "error space id is not a number",
			args: args{
				context: context.Background(),
				spaceID: "invalid",
			},
			want: want{
				err: apperror.NewInvalidData("Invalid space ID format", "strconv.Atoi: parsing \"invalid\": invalid syntax", "space_usecase.go:GetUsersBySpace"),
			},
		},
		{
			name: "error finding space",
			args: args{
				context: context.Background(),
				spaceID: "1",
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error")),
			},
		},
		{
			name: "error finding users",
			args: args{
				context: context.Background(),
				spaceID: "1",
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenSpace, nil),
				mockUserSpaceRepository.EXPECT().FindUserIDsBySpaceID(gomock.Any(), gomock.Any()).Return([]int{}, errors.New("unexpected error")),
			},
		},
		{
			name: "error finding user",
			args: args{
				context: context.Background(),
				spaceID: "1",
			},
			want: want{
				err: errors.New("unexpected error"),
			},
			calls: []*gomock.Call{
				mockSpaceRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(givenSpace, nil),
				mockUserSpaceRepository.EXPECT().FindUserIDsBySpaceID(gomock.Any(), gomock.Any()).Return([]int{1}, nil),
				mockUserRepository.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error")),
			},
		},
	}

	for _, test := range tests {
		calls := make([]interface{}, len(test.calls))
		for i, c := range test.calls {
			calls[i] = c
		}
		gomock.InOrder(calls...)

		got, gotErr := spaceUseCase.GetUsersBySpace(test.args.context, test.args.spaceID)

		assert.Equal(t, test.want.err, gotErr)
		assert.Equal(t, test.want.users, got)
	}
}
