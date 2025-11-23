package notification

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/pkg/helpers"
)

type NotificationUsecase interface {
	SaveNotification(ctx context.Context, notification *domain.Notification) error
	GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]*domain.Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
	MarkAllAsRead(ctx context.Context, userID int) error
	GetUnreadCount(ctx context.Context, userID int) (int, error)
}

type notificationUsecase struct {
	notificationRepo domain.NotificationRepository
}

func NewNotificationUsecase(
	notificationRepo domain.NotificationRepository,
) NotificationUsecase {
	return &notificationUsecase{
		notificationRepo: notificationRepo,
	}
}

func (u *notificationUsecase) SaveNotification(ctx context.Context, notification *domain.Notification) error {
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = helpers.GetTime()
	}

	return u.notificationRepo.SaveNotification(ctx, notification)
}

func (u *notificationUsecase) GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]*domain.Notification, error) {
	return u.notificationRepo.GetUserNotifications(ctx, userID, limit, offset)
}

func (u *notificationUsecase) MarkAsRead(ctx context.Context, notificationID string) error {
	return u.notificationRepo.MarkAsRead(ctx, notificationID)
}

func (u *notificationUsecase) MarkAllAsRead(ctx context.Context, userID int) error {
	return u.notificationRepo.MarkAllAsRead(ctx, userID)
}

func (u *notificationUsecase) GetUnreadCount(ctx context.Context, userID int) (int, error) {
	return u.notificationRepo.GetUnreadCount(ctx, userID)
}
