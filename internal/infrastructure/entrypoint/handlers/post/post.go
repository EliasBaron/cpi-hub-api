package post

import (
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/post"
	"cpi-hub-api/pkg/apperror"
	response "cpi-hub-api/pkg/http"

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

	response.CreatedResponse(c.Writer, "Post created successfully", dto.ToPostWithUserSpaceDTO(createdPost))
}

func (h *PostHandler) Get(c *gin.Context) {
	postID := c.Param("post_id")

	post, err := h.PostUseCase.Get(c.Request.Context(), postID)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, "Post retrieved successfully", dto.ToPostWithUserSpaceDTO(post))
}
