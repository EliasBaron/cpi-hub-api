package events

import (
	"context"
	"cpi-hub-api/internal/core/domain"
)

// EventEmitter es la interfaz para emitir eventos genéricos
type EventEmitter interface {
	EmitEvent(ctx context.Context, event *domain.Event) error
}

type eventEmitter struct {
	notificationManager domain.NotificationManager
}

// NewEventEmitter crea una nueva instancia de EventEmitter
func NewEventEmitter(notificationManager domain.NotificationManager) EventEmitter {
	return &eventEmitter{
		notificationManager: notificationManager,
	}
}

// EmitEvent emite un evento genérico vía WebSocket al usuario target
func (e *eventEmitter) EmitEvent(ctx context.Context, event *domain.Event) error {
	return e.notificationManager.BroadcastEvent(event.TargetUserID, event)
}

