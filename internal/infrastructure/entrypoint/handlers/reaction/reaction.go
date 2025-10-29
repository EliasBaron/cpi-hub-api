package reaction

import (
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/internal/core/usecase/reaction"
	"cpi-hub-api/pkg/apperror"

	response "cpi-hub-api/pkg/http"

	"github.com/gin-gonic/gin"
)

type ReactionHandler struct {
	ReactionUseCase reaction.ReactionUseCase
}

func (h *ReactionHandler) AddReaction(c *gin.Context) {
	var reactionDTO dto.NewReaction
	if err := c.ShouldBindJSON(&reactionDTO); err != nil {
		appErr := apperror.NewInvalidData("Invalid reaction data", err, "reaction_handler.go:AddReaction")
		response.NewError(c.Writer, appErr)
		return
	}

	domainReaction := reactionDTO.ToDomain()
	reaction, err := h.ReactionUseCase.AddReaction(c.Request.Context(), &domainReaction)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.CreatedResponse(c.Writer, dto.ToReactionDTO(*reaction))
}
