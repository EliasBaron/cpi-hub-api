package mapper

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToMongoNotification(notification *domain.Notification) *entity.Notification {
	var oid primitive.ObjectID
	if notification.ID != "" {
		if parsed, err := primitive.ObjectIDFromHex(notification.ID); err == nil {
			oid = parsed
		}
	}

	var createdAt primitive.DateTime
	if !notification.CreatedAt.IsZero() {
		createdAt = primitive.NewDateTimeFromTime(notification.CreatedAt)
	}

	return &entity.Notification{
		ID:          oid,
		Title:       notification.Title,
		Description: notification.Description,
		URL:         notification.URL,
		To:          notification.To,
		Read:        notification.Read,
		CreatedAt:   createdAt,
	}
}

func ToDomainNotification(notificationEntity *entity.Notification) *domain.Notification {
	var idStr string
	if notificationEntity.ID != primitive.NilObjectID {
		idStr = notificationEntity.ID.Hex()
	}

	var createdAt time.Time
	if notificationEntity.CreatedAt != 0 {
		createdAt = notificationEntity.CreatedAt.Time()
	}

	return &domain.Notification{
		ID:          idStr,
		Title:       notificationEntity.Title,
		Description: notificationEntity.Description,
		URL:         notificationEntity.URL,
		To:          notificationEntity.To,
		Read:        notificationEntity.Read,
		CreatedAt:   createdAt,
	}
}
