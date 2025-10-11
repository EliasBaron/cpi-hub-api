package mapper

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/entity"
)

func ToPostgresChatMessage(chatMessage *domain.ChatMessage) *entity.ChatMessageEntity {
	return &entity.ChatMessageEntity{
		ID:        chatMessage.ID,
		Content:   chatMessage.Content,
		UserID:    chatMessage.UserID,
		Username:  chatMessage.Username,
		SpaceID:   chatMessage.SpaceID,
		Timestamp: chatMessage.Timestamp,
	}
}

func ToDomainChatMessage(chatEntity *entity.ChatMessageEntity) *domain.ChatMessage {
	return &domain.ChatMessage{
		ID:        chatEntity.ID,
		Content:   chatEntity.Content,
		UserID:    chatEntity.UserID,
		Username:  chatEntity.Username,
		SpaceID:   chatEntity.SpaceID,
		Timestamp: chatEntity.Timestamp,
	}
}
