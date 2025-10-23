package comment

import (
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/comment"
	"cpi-hub-api/pkg/apperror"
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

func (h *CommentHandler) Update(c *gin.Context) {
	var updateDTO dto.UpdateCommentDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		response.NewError(c.Writer, err)
		return
	}

	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid comment_id (must be integer)", err, "comment_handler.go:Update")
		response.NewError(c.Writer, appErr)
		return
	}

	updateDTO.CommentID = commentID

	err = h.CommentUseCase.Update(c.Request.Context(), updateDTO)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, gin.H{"message": "Comment updated successfully"})
}

func (h *CommentHandler) Delete(c *gin.Context) {
	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid comment_id (must be integer)", err, "comment_handler.go:Delete")
		response.NewError(c.Writer, appErr)
		return
	}

	err = h.CommentUseCase.Delete(c.Request.Context(), commentID)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, gin.H{"message": "Comment deleted successfully"})
}
