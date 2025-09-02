package comment

import (
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/comment"
	"cpi-hub-api/pkg/apperror"
	response "cpi-hub-api/pkg/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	CommentUseCase comment.CommentUseCase
}

func (h *CommentHandler) Create(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid post ID", err, "comment_handler.go:Create")
		response.NewError(c.Writer, appErr)
		return
	}
	var commentDTO dto.CommentDTO

	if err := c.ShouldBindJSON(&commentDTO); err != nil {
		appErr := apperror.NewInvalidData("Invalid comment data", err, "comment_handler.go:Create")
		response.NewError(c.Writer, appErr)
		return
	}
	commentDTO.PostID = postID

	createdComment, err := h.CommentUseCase.Create(c.Request.Context(), commentDTO.ToDomain())

	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.CreatedResponse(c.Writer, "Comment created successfully", dto.ToCommentWithUserAndPostDTO(createdComment))
}

func (h *CommentHandler) Get(c *gin.Context) {
	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid comment ID", err, "comment_handler.go:Get")
		response.NewError(c.Writer, appErr)
		return
	}

	comment, err := h.CommentUseCase.Get(c.Request.Context(), commentID)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, "Comment retrieved successfully", dto.ToCommentWithUserAndPostDTO(comment))
}
