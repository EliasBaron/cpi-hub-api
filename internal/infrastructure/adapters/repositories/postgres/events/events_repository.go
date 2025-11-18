package events

import (
	"cpi-hub-api/internal/core/domain"
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
