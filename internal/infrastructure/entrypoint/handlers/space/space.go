package space

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/space"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
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

	response.CreatedResponse(c.Writer, dto.ToSpaceWithUserDTO(createdSpace))
}

func (h *SpaceHandler) Get(c *gin.Context) {
	spaceId := c.Param("space_id")

	space, err := h.SpaceUseCase.Get(c.Request.Context(), spaceId)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, dto.ToSpaceWithUserDTO(space))
}

func (h *SpaceHandler) Search(context *gin.Context) {
	var searchDTO dto.SearchSpacesDTO

	if err := context.ShouldBindQuery(&searchDTO); err != nil {
		appErr := apperror.NewInvalidData("Invalid query parameters", err, "space_handler.go:Search")
		response.NewError(context.Writer, appErr)
		return
	}

	page, pageSize := helpers.GetPaginationValues(context)
	orderBy, sortDirection := helpers.GetSortValues(context)

	searchCriteria := &domain.SpaceSearchCriteria{
		Name:          searchDTO.Name,
		CreatedBy:     searchDTO.CreatedBy,
		OrderBy:       orderBy,
		Page:          page,
		PageSize:      pageSize,
		SortDirection: sortDirection,
	}

	searchResult, err := h.SpaceUseCase.Search(context.Request.Context(), searchCriteria)
	if err != nil {
		response.NewError(context.Writer, err)
		return
	}

	responseDTO := dto.SearchSpacesResponseDTO{
		Data:     dto.ToSpaceWithUserDTOs(searchResult.Data),
		Page:     searchResult.Page,
		PageSize: searchResult.PageSize,
		Total:    searchResult.Total,
	}

	response.SuccessResponse(context.Writer, responseDTO)
}
