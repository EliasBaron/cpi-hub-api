package events

import (
	"cpi-hub-api/internal/core/dto"
	eventsUsecase "cpi-hub-api/internal/core/usecase/events"
	"cpi-hub-api/pkg/apperror"
	response "cpi-hub-api/pkg/http"

	"github.com/gin-gonic/gin"
)

// EventsHandler maneja las conexiones de eventos en tiempo real
type EventsHandler struct {
	eventsUsecase *eventsUsecase.EventsUsecase
}

// NewEventsHandler crea una nueva instancia del handler
func NewEventsHandler(eventsUsecase *eventsUsecase.EventsUsecase) *EventsHandler {
	return &EventsHandler{
		eventsUsecase: eventsUsecase,
	}
}

// Connect maneja la conexión de eventos en tiempo real
func (h *EventsHandler) Connect(c *gin.Context) {
	spaceID := c.Param("space_id")
	if spaceID == "" {
		appErr := apperror.NewInvalidData("space_id es requerido", nil, "events_handler.go:Connect")
		response.NewError(c.Writer, appErr)
		return
	}

	userID := c.Query("user_id")
	if userID == "" {
		appErr := apperror.NewInvalidData("user_id es requerido", nil, "events_handler.go:Connect")
		response.NewError(c.Writer, appErr)
		return
	}

	connectionParams := dto.EventsConnectionParams{
		UserID:  userID,
		SpaceID: spaceID,
		Writer:  c.Writer,
		Request: c.Request,
	}

	err := h.eventsUsecase.HandleConnection(connectionParams)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}
}

// Broadcast envía un mensaje a todos los usuarios
func (h *EventsHandler) Broadcast(c *gin.Context) {
	spaceID := c.Param("space_id")
	if spaceID == "" {
		appErr := apperror.NewInvalidData("space_id es requerido", nil, "events_handler.go:Broadcast")
		response.NewError(c.Writer, appErr)
		return
	}

	var dto dto.EventsBroadcastParams
	if err := c.ShouldBindJSON(&dto); err != nil {
		appErr := apperror.NewInvalidData("Invalid request data", err, "events_handler.go:Broadcast")
		response.NewError(c.Writer, appErr)
		return
	}

	dto.SpaceID = spaceID

	chatMsg, err := h.eventsUsecase.Broadcast(dto)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, gin.H{
		"message":      "Mensaje enviado exitosamente",
		"chat_message": chatMsg,
	})
}
