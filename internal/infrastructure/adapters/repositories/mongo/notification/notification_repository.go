package notification

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/mapper"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationRepository struct {
	db *mongo.Database
}

func NewNotificationRepository(db *mongo.Database) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) SaveNotification(ctx context.Context, notification *domain.Notification) error {
	notificationEntity := mapper.ToMongoNotification(notification)

	collection := r.db.Collection("notifications")

	res, err := collection.InsertOne(ctx, notificationEntity)
	if err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		notification.ID = oid.Hex()
	}

	return nil
}

func (r *NotificationRepository) GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]*domain.Notification, error) {
	collection := r.db.Collection("notifications")

	filter := bson.M{"user_id": userID}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}
	defer cursor.Close(ctx)

	var notifications []*domain.Notification
	for cursor.Next(ctx) {
		var notificationEntity entity.Notification

		if err := cursor.Decode(&notificationEntity); err != nil {
			return nil, fmt.Errorf("failed to decode notification: %w", err)
		}

		notification := mapper.ToDomainNotification(&notificationEntity)
		notifications = append(notifications, notification)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return notifications, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID string) error {
	oid, err := primitive.ObjectIDFromHex(notificationID)
	if err != nil {
		return fmt.Errorf("invalid notification ID: %w", err)
	}

	collection := r.db.Collection("notifications")
	update := bson.M{
		"$set": bson.M{"read": true},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID int) error {
	collection := r.db.Collection("notifications")
	filter := bson.M{"user_id": userID, "read": false}
	update := bson.M{
		"$set": bson.M{"read": true},
	}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID int) (int, error) {
	collection := r.db.Collection("notifications")
	filter := bson.M{"user_id": userID, "read": false}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return int(count), nil
}
