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
