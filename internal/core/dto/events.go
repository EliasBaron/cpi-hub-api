package dto

import (
	"cpi-hub-api/internal/core/domain"
	"net/http"
	"time"
)

type EventsConnectionParams struct {
	UserID   int
	SpaceID  int
	Writer   http.ResponseWriter
	Request  *http.Request
	Username string
}

type EventsBroadcastParams struct {
	SpaceID  int    `json:"space_id" binding:"required"`
	UserID   int    `json:"user_id" binding:"required"`
	Message  string `json:"message" binding:"required"`
	Username string `json:"username" binding:"required"`
	Image    string `json:"image"`
}

type HandleUserConnectionParams struct {
	UserID  int
	Writer  http.ResponseWriter
	Request *http.Request
}

type NotificationMessageDTO struct {
	Type      string          `json:"type"`
	Data      NotificationDTO `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
}

func ToNotificationMessageDTO(notification *domain.Notification) NotificationMessageDTO {
	return NotificationMessageDTO{
		Type:      "notification",
		Data:      ToNotificationDTO(notification),
		Timestamp: notification.CreatedAt,
	}
}

type EventMessageDTO struct {
	Type      string        `json:"type"`
	Data      *domain.Event `json:"data"`
	Timestamp time.Time     `json:"timestamp"`
}

func ToEventMessageDTO(event *domain.Event) EventMessageDTO {
	return EventMessageDTO{
		Type:      "event",
		Data:      event,
		Timestamp: event.Timestamp,
	}
}
