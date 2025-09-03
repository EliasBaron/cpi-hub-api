package post

import (
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/post"
	"cpi-hub-api/pkg/apperror"
	response "cpi-hub-api/pkg/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	PostUseCase post.PostUseCase
}

func (h *PostHandler) Create(c *gin.Context) {
	var postDTO dto.CreatePost

	if err := c.ShouldBindJSON(&postDTO); err != nil {
		appErr := apperror.NewInvalidData("Invalid post data", err, "post_handler.go:Create")
		response.NewError(c.Writer, appErr)
		return
	}

	createdPost, err := h.PostUseCase.Create(c.Request.Context(), postDTO.ToDomain())
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.CreatedResponse(c.Writer, "Post created successfully", dto.ToPostExtendedDTO(createdPost))
}

func (h *PostHandler) Get(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid post ID", err, "post_handler.go:Get")
		response.NewError(c.Writer, appErr)
		return
	}

	post, err := h.PostUseCase.Get(c.Request.Context(), postID)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, "Post retrieved successfully", dto.ToPostExtendedDTO(post))
}

func (h *PostHandler) AddComment(c *gin.Context) {
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

	createdComment, err := h.PostUseCase.AddComent(c.Request.Context(), commentDTO.ToDomain())

	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.CreatedResponse(c.Writer, "Comment created successfully", dto.ToCommentWithUserAndPostDTO(createdComment))
}
