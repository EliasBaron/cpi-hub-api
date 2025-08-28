package user

import (
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/user"
	"cpi-hub-api/pkg/apperror"
	response "cpi-hub-api/pkg/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UseCase user.UserUseCase
}

func (h *Handler) Create(c *gin.Context) {
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

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")

	user, err := h.UseCase.Get(c.Request.Context(), id)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, "User retrieved successfully", dto.ToUserDTOWithSpaces(user))
}

func (h *Handler) AddSpaceToUser(c *gin.Context) {
	userId := c.Param("id")
	spaceId := c.Param("spaceId")

	err := h.UseCase.AddSpaceToUser(c.Request.Context(), userId, spaceId)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, "Space added to user successfully", nil)
}
