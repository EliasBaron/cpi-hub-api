package post

import (
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/post"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
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

	createdComment, err := h.PostUseCase.AddComment(c.Request.Context(), commentDTO.ToDomain())

	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.CreatedResponse(c.Writer, "Comment created successfully", dto.ToCommentWithUserAndPostDTO(createdComment))
}

func (h *PostHandler) Search(context *gin.Context) {
	page, pageSize := helpers.GetPaginationValues(context)
	orderBy, sortDirection := helpers.GetSortValues(context)

	var spaceID int
	if spaceIDStr := context.Query("space_id"); spaceIDStr != "" {
		if id, err := strconv.Atoi(spaceIDStr); err == nil && id > 0 {
			spaceID = id
		} else {
			appErr := apperror.NewInvalidData("Invalid space_id parameter (must be positive integer)", err, "post_handler.go:Search")
			response.NewError(context.Writer, appErr)
			return
		}
	}

	searchParams := dto.SearchPostsParams{
		Page:          page,
		PageSize:      pageSize,
		OrderBy:       orderBy,
		SortDirection: sortDirection,
		SpaceID:       spaceID,
		Query:         context.Query("q"),
	}

	searchResult, err := h.PostUseCase.Search(context.Request.Context(), searchParams)
	if err != nil {
		response.NewError(context.Writer, err)
		return
	}

	data := dto.PaginatedPostsResponse{
		Data:     dto.ToPostExtendedDTOs(searchResult.Posts),
		Page:     searchParams.Page,
		PageSize: searchParams.PageSize,
		Total:    searchResult.Total,
	}

	response.SuccessResponse(context.Writer, "Posts retrieved successfully", data)
}

func (h *PostHandler) GetInterestedPosts(context *gin.Context) {
	page, pageSize := helpers.GetPaginationValues(context)
	orderBy, sortDirection := helpers.GetSortValues(context)

	userIDStr := context.Query("user_id")
	if userIDStr == "" {
		appErr := apperror.NewInvalidData("user_id parameter is required", nil, "post_handler.go:GetInterestedPosts")
		response.NewError(context.Writer, appErr)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		appErr := apperror.NewInvalidData("Invalid user_id parameter (must be positive integer)", err, "post_handler.go:GetInterestedPosts")
		response.NewError(context.Writer, appErr)
		return
	}

	interestedParams := dto.InterestedPostsParams{
		Page:          page,
		PageSize:      pageSize,
		OrderBy:       orderBy,
		SortDirection: sortDirection,
		UserID:        userID,
	}

	searchResult, err := h.PostUseCase.GetInterestedPosts(context.Request.Context(), interestedParams)
	if err != nil {
		response.NewError(context.Writer, err)
		return
	}

	data := dto.PaginatedPostsResponse{
		Data:     dto.ToPostExtendedDTOs(searchResult.Posts),
		Page:     interestedParams.Page,
		PageSize: interestedParams.PageSize,
		Total:    searchResult.Total,
	}

	response.SuccessResponse(context.Writer, "Interested posts retrieved successfully", data)
}
