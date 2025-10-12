package message

import (
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/message"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
	response "cpi-hub-api/pkg/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	MessageUseCase message.MessageUseCase
}

func (h *MessageHandler) Search(c *gin.Context) {
	page, pageSize := helpers.GetPaginationValues(c)
	orderBy, sortDirection := helpers.GetSortValues(c)

	var spaceID int
	if spaceIDStr := c.Query("space_id"); spaceIDStr != "" {
		if id, err := strconv.Atoi(spaceIDStr); err == nil && id > 0 {
			spaceID = id
		} else {
			appErr := apperror.NewInvalidData("Invalid space_id parameter (must be positive integer)", err, "message_handler.go:Search")
			response.NewError(c.Writer, appErr)
			return
		}
	} else {
		appErr := apperror.NewInvalidData("space_id parameter is required", nil, "message_handler.go:Search")
		response.NewError(c.Writer, appErr)
		return
	}

	searchParams := dto.SearchMessagesParams{
		Page:          page,
		PageSize:      pageSize,
		OrderBy:       orderBy,
		SortDirection: sortDirection,
		SpaceID:       spaceID,
	}

	searchResult, err := h.MessageUseCase.Search(c.Request.Context(), searchParams)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	data := dto.PaginatedMessagesResponse{
		Data:     dto.ToMessageDTOs(searchResult.Messages),
		Page:     searchParams.Page,
		PageSize: searchParams.PageSize,
		Total:    searchResult.Total,
	}

	response.SuccessResponse(c.Writer, data)
}
