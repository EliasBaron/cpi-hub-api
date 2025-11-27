package user

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/post"
	"cpi-hub-api/internal/core/usecase/user"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
	response "cpi-hub-api/pkg/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UseCase     user.UserUseCase
	PostUseCase post.PostUseCase
}

func (h *UserHandler) Register(c *gin.Context) {
	var createUserDTO dto.CreateUser

	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		appErr := apperror.NewInvalidData("Invalid user data", err, "user_handler.go:Create")
		response.NewError(c.Writer, appErr)
		return
	}

	user := createUserDTO.ToDomain()

	createdUser, err := h.UseCase.Create(c.Request.Context(), user)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	token, err := helpers.CreateToken(createdUser.Email, createdUser.ID)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	c.Header("Authorization", "Bearer "+token)

	registerResponse := dto.AuthResponse{
		User:  dto.ToUserDTO(createdUser),
		Token: token,
	}

	response.CreatedResponse(c.Writer, registerResponse)
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginDTO dto.LoginUser

	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		appErr := apperror.NewInvalidData("Invalid login data", err, "user_handler.go:Login")
		response.NewError(c.Writer, appErr)
		return
	}

	user, err := h.UseCase.Login(c.Request.Context(), loginDTO)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	token, err := helpers.CreateToken(user.Email, user.ID)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	c.Header("Authorization", "Bearer "+token)

	loginResponse := dto.AuthResponse{
		User:  dto.ToUserDTO(user),
		Token: token,
	}

	response.SuccessResponse(c.Writer, loginResponse)
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		response.NewError(c.Writer, apperror.NewUnauthorized("Missing authorization token", nil, "user_handler.go:GetCurrentUser"))
		return
	}

	userId, err := helpers.GetUserIdFromToken(token)
	if err != nil {
		response.NewError(c.Writer, apperror.NewUnauthorized("Invalid token", err, "user_handler.go:GetCurrentUser"))
		return
	}

	user, err := h.UseCase.Get(c.Request.Context(), userId)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, dto.ToUserDTOWithSpaces(user))
}

func (h *UserHandler) Get(c *gin.Context) {
	idStr := c.Param("user_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid user_id (must be integer)", err, "user_handler.go:Get")
		response.NewError(c.Writer, appErr)
		return
	}

	user, err := h.UseCase.Get(c.Request.Context(), id)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, dto.ToUserDTOWithSpaces(user))
}

func (h *UserHandler) AddSpaceToUser(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid user_id (must be integer)", err, "user_handler.go:AddSpaceToUser")
		response.NewError(c.Writer, appErr)
		return
	}

	spaceIdStr := c.Param("space_id")
	spaceId, err := strconv.Atoi(spaceIdStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid space_id (must be integer)", err, "user_handler.go:AddSpaceToUser")
		response.NewError(c.Writer, appErr)
		return
	}

	err = h.UseCase.Update(c.Request.Context(), dto.UpdateUserSpacesDTO{
		UserID:   userId,
		SpaceIDs: []int{spaceId},
		Action:   domain.AddUserToSpace,
	})

	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, nil)
}

func (h *UserHandler) RemoveSpaceFromUser(c *gin.Context) {
	userIdStr, spaceIdStr := c.Param("user_id"), c.Param("space_id")
	userId, err1 := strconv.Atoi(userIdStr)
	spaceId, err2 := strconv.Atoi(spaceIdStr)

	if err1 != nil || err2 != nil {
		appErr := apperror.NewInvalidData("Invalid space_id or user_id (must be integer)", nil, "user_handler.go:RemoveSpaceFromUser")
		response.NewError(c.Writer, appErr)
		return
	}

	err := h.UseCase.Update(c.Request.Context(), dto.UpdateUserSpacesDTO{
		UserID:   userId,
		SpaceIDs: []int{spaceId},
		Action:   domain.RemoveUserFromSpace,
	})
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, nil)
}

func (h *UserHandler) GetSpacesByUserId(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid user_id (must be integer)", err, "user_handler.go:GetSpacesByUserId")
		response.NewError(c.Writer, appErr)
		return
	}

	spaces, err := h.UseCase.GetSpacesByUser(c.Request.Context(), userId)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	spacesDTO := make([]dto.SpaceDTO, len(spaces))
	for i, space := range spaces {
		spacesDTO[i] = dto.ToSpaceDTO(space)
	}

	response.SuccessResponse(c.Writer, spacesDTO)
}

func (h *UserHandler) GetTrendingUsers(c *gin.Context) {
	page, pageSize := helpers.GetPaginationValues(c)
	timeFrame := c.Query("time_frame")

	trendingParams := dto.TrendingUsersParams{
		Page:      page,
		PageSize:  pageSize,
		TimeFrame: timeFrame,
	}

	trendingUsers, err := h.UseCase.GetTrendingUsers(c.Request.Context(), trendingParams)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, trendingUsers)
}

func (h *UserHandler) GetInterestedPosts(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid user_id (must be integer)", err, "user_handler.go:GetInterestedPosts")
		response.NewError(c.Writer, appErr)
		return
	}

	page, pageSize := helpers.GetPaginationValues(c)
	orderBy, sortDirection := helpers.GetSortValues(c)

	searchResult, err := h.PostUseCase.GetInterestedPosts(c.Request.Context(), dto.InterestedPostsParams{
		Page:          page,
		PageSize:      pageSize,
		OrderBy:       orderBy,
		SortDirection: sortDirection,
		UserID:        userId,
	})
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	data := dto.PaginatedPostsResponse{
		Data:     dto.ToPostExtendedDTOs(searchResult.Posts),
		Page:     page,
		PageSize: pageSize,
		Total:    searchResult.Total,
	}

	response.SuccessResponse(c.Writer, data)
}

func (h *UserHandler) Search(c *gin.Context) {
	var searchParams dto.SearchUsersParams

	if err := c.ShouldBindQuery(&searchParams); err != nil {
		appErr := apperror.NewInvalidData("Invalid search parameters", err, "user_handler.go:Search")
		response.NewError(c.Writer, appErr)
		return
	}

	page, pageSize := helpers.GetPaginationValues(c)
	orderBy, sortDirection := helpers.GetSortValues(c)

	searchParams.Page = page
	searchParams.PageSize = pageSize
	searchParams.OrderBy = orderBy
	searchParams.SortDirection = sortDirection

	result, err := h.UseCase.Search(c.Request.Context(), searchParams)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, result)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var updateUserDTO dto.UpdateUserDTO

	if err := c.ShouldBindJSON(&updateUserDTO); err != nil {
		appErr := apperror.NewInvalidData("Invalid user data", err, "user_handler.go:UpdateUser")
		response.NewError(c.Writer, appErr)
		return
	}

	userIdStr := c.Param("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid user_id (must be integer)", err, "user_handler.go:UpdateUser")
		response.NewError(c.Writer, appErr)
		return
	}

	updateUserDTO.UserID = userId

	err = h.UseCase.UpdateUser(c.Request.Context(), updateUserDTO)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, nil)
}
