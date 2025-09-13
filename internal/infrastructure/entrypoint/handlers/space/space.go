package space

import (
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/space"
	"cpi-hub-api/pkg/apperror"
	response "cpi-hub-api/pkg/http"

	"github.com/gin-gonic/gin"
)

type SpaceHandler struct {
	SpaceUseCase space.SpaceUseCase
}

func (h *SpaceHandler) Create(c *gin.Context) {
	var createSpaceDTO dto.CreateSpace

	if err := c.ShouldBindJSON(&createSpaceDTO); err != nil {
		appErr := apperror.NewInvalidData("Invalid space data", err, "space_handler.go:Create")
		response.NewError(c.Writer, appErr)
		return
	}

	createdSpace, err := h.SpaceUseCase.Create(c.Request.Context(), createSpaceDTO.ToDomain())
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.CreatedResponse(c.Writer, "Space created successfully", dto.ToSpaceWithUserDTO(createdSpace))
}

func (h *SpaceHandler) Get(c *gin.Context) {
	spaceId := c.Param("space_id")

	space, err := h.SpaceUseCase.Get(c.Request.Context(), spaceId)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, "Space retrieved successfully", dto.ToSpaceWithUserDTO(space))
}

func (h *SpaceHandler) Search(context *gin.Context) {
	orderBy := context.Query("order_by")
	if orderBy == "" {
		appErr := apperror.NewInvalidData("order_by query parameter is required", nil, "space_handler.go:Search")
		response.NewError(context.Writer, appErr)
		return
	}
	spaces, err := h.SpaceUseCase.GetAll(context.Request.Context(), orderBy)
	if err != nil {
		response.NewError(context.Writer, err)
		return
	}

	response.SuccessResponse(context.Writer, "Spaces retrieved successfully", dto.ToSpaceWithUserDTOs(spaces))
}
