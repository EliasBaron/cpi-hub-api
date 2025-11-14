package dto

import (
	"cpi-hub-api/internal/core/domain"
	"net/http"
	"time"
)

type HandleNotificationConnectionParams struct {
	UserID  int
	Writer  http.ResponseWriter
	Request *http.Request
}

type CreateNotificationParams struct {
	NotificationType domain.NotificationType
	EntityType       domain.EntityType
	EntityID         int
	PostID           *int
	OwnerUserID      int
}

type NotificationDTO struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	EntityType string    `json:"entity_type"`
	EntityID   int       `json:"entity_id"`
	PostID     *int      `json:"post_id,omitempty"` // PostID is set when EntityType is comment
	UserID     int       `json:"user_id"`
	Read       bool      `json:"read"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToNotificationDTO(notification *domain.Notification) NotificationDTO {
	return NotificationDTO{
		ID:         notification.ID,
		Type:       string(notification.Type),
		EntityType: string(notification.EntityType),
		EntityID:   notification.EntityID,
		PostID:     notification.PostID,
		UserID:     notification.UserID,
		Read:       notification.Read,
		CreatedAt:  notification.CreatedAt,
	}
}
