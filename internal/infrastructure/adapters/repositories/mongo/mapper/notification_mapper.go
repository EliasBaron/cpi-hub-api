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
		ID:         oid,
		Type:       string(notification.Type),
		EntityType: string(notification.EntityType),
		EntityID:   notification.EntityID,
		UserID:     notification.UserID,
		Read:       notification.Read,
		CreatedAt:  createdAt,
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
		ID:         idStr,
		Type:       domain.NotificationType(notificationEntity.Type),
		EntityType: domain.EntityType(notificationEntity.EntityType),
		EntityID:   notificationEntity.EntityID,
		UserID:     notificationEntity.UserID,
		Read:       notificationEntity.Read,
		CreatedAt:  createdAt,
	}
}
