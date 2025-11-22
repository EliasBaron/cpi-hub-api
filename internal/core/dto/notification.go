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

// CreateNotificationDTO representa una notificación genérica para crear desde el frontend
type CreateNotificationDTO struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	URL         *string `json:"url,omitempty"`
	To          int     `json:"to" binding:"required"`
}

// ToDomain convierte CreateNotificationDTO a domain.Notification
func (dto *CreateNotificationDTO) ToDomain() *domain.Notification {
	return &domain.Notification{
		Title:       dto.Title,
		Description: dto.Description,
		URL:         dto.URL,
		To:          dto.To,
		Read:        false,
		CreatedAt:   time.Now(),
	}
}

// NotificationDTO representa una notificación genérica
type NotificationDTO struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         *string   `json:"url,omitempty"`
	To          int       `json:"to"`
	Read        bool      `json:"read"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToNotificationDTO convierte domain.Notification a NotificationDTO
func ToNotificationDTO(notification *domain.Notification) NotificationDTO {
	return NotificationDTO{
		ID:          notification.ID,
		Title:       notification.Title,
		Description: notification.Description,
		URL:         notification.URL,
		To:          notification.To,
		Read:        notification.Read,
		CreatedAt:   notification.CreatedAt,
	}
}
