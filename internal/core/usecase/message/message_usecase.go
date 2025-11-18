package message

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/dto"
)

type SearchResult struct {
	Messages []*domain.ChatMessage
	Total    int
}

type MessageUseCase interface {
	Search(ctx context.Context, params dto.SearchMessagesParams) (*SearchResult, error)
}

type messageUseCase struct {
	messageRepository domain.MessageRepository
}

func NewMessageUsecase(messageRepo domain.MessageRepository) MessageUseCase {
	return &messageUseCase{
		messageRepository: messageRepo,
	}
}

func (m *messageUseCase) Search(ctx context.Context, params dto.SearchMessagesParams) (*SearchResult, error) {
	filters := domain.SearchMessagesFilter{
		SpaceID:       params.SpaceID,
		Page:          params.Page,
		PageSize:      params.PageSize,
		OrderBy:       params.OrderBy,
		SortDirection: params.SortDirection,
	}

	messages, total, err := m.messageRepository.SearchMessages(
		ctx,
		filters,
	)
	if err != nil {
		return nil, err
	}

	return &SearchResult{
		Messages: messages,
		Total:    total,
	}, nil
}
