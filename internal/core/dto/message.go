package dto

import (
	"cpi-hub-api/internal/core/domain"
	"time"
)

type SearchMessagesParams struct {
	Page          int
	PageSize      int
	OrderBy       string
	SortDirection string
	SpaceID       int
}

type PaginatedMessagesResponse struct {
	Data     []MessageDTO `json:"data"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
	Total    int          `json:"total"`
}

type MessageDTO struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	SpaceID   int       `json:"space_id"`
	CreatedAt time.Time `json:"created_at"`
	Image     string    `json:"image"`
}

func ToMessageDTO(message *domain.ChatMessage) MessageDTO {
	return MessageDTO{
		ID:        message.ID,
		Content:   message.Content,
		UserID:    message.UserID,
		Username:  message.Username,
		SpaceID:   message.SpaceID,
		CreatedAt: message.Timestamp,
		Image:     message.Image,
	}
}

func ToMessageDTOs(messages []*domain.ChatMessage) []MessageDTO {
	messageDTOs := make([]MessageDTO, 0, len(messages))

	for _, message := range messages {
		messageDTOs = append(messageDTOs, ToMessageDTO(message))
	}

	return messageDTOs
}
