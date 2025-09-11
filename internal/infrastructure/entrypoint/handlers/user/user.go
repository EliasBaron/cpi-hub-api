package user

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/post"
	"cpi-hub-api/internal/core/usecase/user"
	"cpi-hub-api/pkg/apperror"
	response "cpi-hub-api/pkg/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UseCase     user.UserUseCase
	PostUseCase post.PostUseCase
}

func (h *UserHandler) Create(c *gin.Context) {
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

	response.CreatedResponse(c.Writer, "User created successfully", dto.ToUserDTO(createdUser))
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

	response.SuccessResponse(c.Writer, "User retrieved successfully", dto.ToUserDTOWithSpaces(user))
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

	response.SuccessResponse(c.Writer, "Space added to user successfully", nil)
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

	response.SuccessResponse(c.Writer, "Space removed from user successfully", nil)
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

	response.SuccessResponse(c.Writer, "Spaces retrieved successfully", spacesDTO)
}

func (h *UserHandler) GetPostsByUserSpaces(c *gin.Context) {
	userIDstr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid user ID", err, "user_handler.go:GetPostsByUserSpaces")
		response.NewError(c.Writer, appErr)
		return
	}

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	posts, err := h.PostUseCase.GetPostsByUserSpaces(c.Request.Context(), userID, page)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, "Posts retrieved successfully", dto.ToPostExtendedDTOs(posts))
}
