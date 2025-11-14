package notification

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/pkg/helpers"
	"log"
)

type NotificationUsecase interface {
	CreateNotification(ctx context.Context, params dto.CreateNotificationParams) error
	GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]*domain.Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
	MarkAllAsRead(ctx context.Context, userID int) error
	GetUnreadCount(ctx context.Context, userID int) (int, error)
}

type notificationUsecase struct {
	notificationRepo domain.NotificationRepository
	notificationMgr  domain.NotificationManager
}

func NewNotificationUsecase(
	notificationRepo domain.NotificationRepository,
	notificationMgr domain.NotificationManager,
) NotificationUsecase {
	return &notificationUsecase{
		notificationRepo: notificationRepo,
		notificationMgr:  notificationMgr,
	}
}

func (u *notificationUsecase) CreateNotification(ctx context.Context, params dto.CreateNotificationParams) error {
	notification := &domain.Notification{
		Type:       params.NotificationType,
		EntityType: params.EntityType,
		EntityID:   params.EntityID,
		UserID:     params.OwnerUserID,
		Read:       false,
		CreatedAt:  helpers.GetTime(),
	}

	err := u.notificationRepo.SaveNotification(ctx, notification)
	if err != nil {
		return err
	}

	err = u.notificationMgr.BroadcastToUser(params.OwnerUserID, notification)
	if err != nil {
		log.Printf("Error broadcasting notification to user %d: %v", params.OwnerUserID, err)
	}

	return nil
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
