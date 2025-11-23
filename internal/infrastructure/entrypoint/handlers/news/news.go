package news

import (
	"cpi-hub-api/internal/core/dto"
	newsUsecase "cpi-hub-api/internal/core/usecase/news"

	response "cpi-hub-api/pkg/http"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	NewsUseCase newsUsecase.NewsUseCase
}

func NewNewsHandler(newsUseCase newsUsecase.NewsUseCase) *NewsHandler {
	return &NewsHandler{
		NewsUseCase: newsUseCase,
	}
}

func (h *NewsHandler) GetAll(c *gin.Context) {

	newsItems, err := h.NewsUseCase.GetAllNews(c.Request.Context())
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, newsItems)
}

func (h *NewsHandler) Create(c *gin.Context) {

	var createNewsDTO dto.CreateNewsDTO
	if err := c.ShouldBindJSON(&createNewsDTO); err != nil {
		response.NewError(c.Writer, err)
		return
	}

	newsItem, err := h.NewsUseCase.CreateNews(c.Request.Context(), createNewsDTO.ToDomain())
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, newsItem)
}
