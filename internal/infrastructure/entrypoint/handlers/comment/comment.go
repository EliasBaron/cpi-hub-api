package comment

import (
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/comment"
	"cpi-hub-api/pkg/helpers"
	response "cpi-hub-api/pkg/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	CommentUseCase comment.CommentUseCase
}

func (h *CommentHandler) Search(c *gin.Context) {
	page, pageSize := helpers.GetPaginationValues(c)
	orderBy, sortDirection := helpers.GetSortValues(c)

	var userIDInt int
	var postIDInt int

	if userID := c.Query("user_id"); userID != "" {
		userIDInt, _ = strconv.Atoi(userID)
	}

	if postID := c.Query("post_id"); postID != "" {
		postIDInt, _ = strconv.Atoi(postID)
	}

	searchParams := dto.SearchCommentsParams{
		Page:          page,
		PageSize:      pageSize,
		OrderBy:       orderBy,
		SortDirection: sortDirection,
		UserID:        &userIDInt,
		PostID:        &postIDInt,
	}

	searchResult, err := h.CommentUseCase.Search(c.Request.Context(), searchParams)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	data := dto.PaginatedCommentsResponse{
		Data:     dto.ToCommentWithSpaceDTOs(searchResult.Comments),
		Page:     searchParams.Page,
		PageSize: searchParams.PageSize,
		Total:    searchResult.Total,
	}

	response.SuccessResponse(c.Writer, data)
}
