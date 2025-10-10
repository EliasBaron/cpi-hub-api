package events

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/mapper"
	"database/sql"
)

type EventsRepository struct {
	db *sql.DB
}

func NewEventsRepository(db *sql.DB) *EventsRepository {
	return &EventsRepository{
		db: db,
	}
}

func (r *EventsRepository) SaveMessage(message *domain.ChatMessage) error {
	chatEntity := mapper.ToPostgresChatMessage(message)

	query := `
		INSERT INTO chat_messages (id, content, user_id, username, space_id, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query,
		chatEntity.ID,
		chatEntity.Content,
		chatEntity.UserID,
		chatEntity.Username,
		chatEntity.SpaceID,
		chatEntity.Timestamp,
	)

	return err
}

func (r *EventsRepository) GetMessagesBySpace(spaceID string, limit int) ([]*domain.ChatMessage, error) {
	query := `
		SELECT id, content, user_id, username, space_id, timestamp
		FROM chat_messages
		WHERE space_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, spaceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.ChatMessage
	for rows.Next() {
		var chatEntity entity.ChatMessageEntity
		err := rows.Scan(
			&chatEntity.ID,
			&chatEntity.Content,
			&chatEntity.UserID,
			&chatEntity.Username,
			&chatEntity.SpaceID,
			&chatEntity.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		message := mapper.ToDomainChatMessage(&chatEntity)
		messages = append(messages, message)
	}

	return messages, nil
}
