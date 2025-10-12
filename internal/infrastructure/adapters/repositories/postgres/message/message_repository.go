package message

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/mapper"
	"database/sql"
	"fmt"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) SearchMessages(ctx context.Context, filters domain.SearchMessagesFilter) ([]*domain.ChatMessage, int, error) {
	baseQuery := "SELECT m.id, m.content, m.user_id, m.username, m.space_id, m.timestamp, u.image FROM chat_messages m INNER JOIN users u ON m.user_id = u.id WHERE m.space_id = $1"
	countQuery := "SELECT COUNT(*) FROM chat_messages WHERE space_id = $1"

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, filters.SpaceID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	orderClause := "ORDER BY timestamp ASC"
	if filters.OrderBy != "" {
		validOrderFields := map[string]bool{
			"timestamp": true,
			"id":        true,
		}

		if validOrderFields[filters.OrderBy] {
			direction := "ASC"
			if filters.SortDirection == "desc" {
				direction = "DESC"
			}
			orderClause = fmt.Sprintf("ORDER BY %s %s", filters.OrderBy, direction)
		}
	}

	offset := (filters.Page - 1) * filters.PageSize
	limitClause := fmt.Sprintf("LIMIT %d OFFSET %d", filters.PageSize, offset)

	query := fmt.Sprintf("%s %s %s", baseQuery, orderClause, limitClause)

	rows, err := r.db.QueryContext(ctx, query, filters.SpaceID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []*domain.ChatMessage
	for rows.Next() {
		var chatEntity entity.ChatMessageEntityWithUser
		err := rows.Scan(
			&chatEntity.ID,
			&chatEntity.Content,
			&chatEntity.UserID,
			&chatEntity.Username,
			&chatEntity.SpaceID,
			&chatEntity.Timestamp,
			&chatEntity.Image,
		)
		if err != nil {
			return nil, 0, err
		}

		message := mapper.ToDomainChatMessageWithUser(&chatEntity)
		messages = append(messages, message)
	}

	return messages, total, rows.Err()
}
