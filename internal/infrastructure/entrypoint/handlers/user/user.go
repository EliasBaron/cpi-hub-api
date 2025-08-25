package user

import (
	"cpi-hub-api/internal/core/usecase/user"
	response "cpi-hub-api/pkg/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UseCase user.UseCase
}

func (h *Handler) Create(c *gin.Context) {
	response.CreatedResponse(c.Writer, "User created successfully", nil)
}
