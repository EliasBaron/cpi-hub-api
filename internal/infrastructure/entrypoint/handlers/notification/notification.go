package notification

import (
	"cpi-hub-api/internal/core/dto"
	notificationUsecase "cpi-hub-api/internal/core/usecase/notification"
	"cpi-hub-api/pkg/apperror"
	response "cpi-hub-api/pkg/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	NotificationUseCase notificationUsecase.NotificationUsecase
}

func NewNotificationHandler(notificationUseCase notificationUsecase.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{
		NotificationUseCase: notificationUseCase,
	}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		appErr := apperror.NewInvalidData("user_id is required", nil, "notification_handler.go:GetNotifications")
		response.NewError(c.Writer, appErr)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid user_id (must be integer)", err, "notification_handler.go:GetNotifications")
		response.NewError(c.Writer, appErr)
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	notifications, err := h.NotificationUseCase.GetUserNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	notificationDTOs := make([]dto.NotificationDTO, len(notifications))
	for i, notif := range notifications {
		notificationDTOs[i] = dto.ToNotificationDTO(notif)
	}

	response.SuccessResponse(c.Writer, gin.H{
		"data":   notificationDTOs,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		appErr := apperror.NewInvalidData("user_id is required", nil, "notification_handler.go:GetUnreadCount")
		response.NewError(c.Writer, appErr)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid user_id (must be integer)", err, "notification_handler.go:GetUnreadCount")
		response.NewError(c.Writer, appErr)
		return
	}

	count, err := h.NotificationUseCase.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, gin.H{
		"unread_count": count,
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationID := c.Param("notification_id")
	if notificationID == "" {
		appErr := apperror.NewInvalidData("notification_id is required", nil, "notification_handler.go:MarkAsRead")
		response.NewError(c.Writer, appErr)
		return
	}

	err := h.NotificationUseCase.MarkAsRead(c.Request.Context(), notificationID)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, gin.H{
		"message": "Notification marked as read",
	})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		appErr := apperror.NewInvalidData("user_id is required", nil, "notification_handler.go:MarkAllAsRead")
		response.NewError(c.Writer, appErr)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		appErr := apperror.NewInvalidData("Invalid user_id (must be integer)", err, "notification_handler.go:MarkAllAsRead")
		response.NewError(c.Writer, appErr)
		return
	}

	err = h.NotificationUseCase.MarkAllAsRead(c.Request.Context(), userID)
	if err != nil {
		response.NewError(c.Writer, err)
		return
	}

	response.SuccessResponse(c.Writer, gin.H{
		"message": "All notifications marked as read",
	})
}
