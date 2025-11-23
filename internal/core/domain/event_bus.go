package domain

import "context"

// EventBus es la interfaz para publicar y suscribirse a eventos del dominio
type EventBus interface {
	Publish(ctx context.Context, event DomainEvent) error
	Subscribe(eventType string, handler EventHandler)
}

// EventHandler es la interfaz para manejar eventos del dominio
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
}
