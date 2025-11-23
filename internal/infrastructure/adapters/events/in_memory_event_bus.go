package events

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"sync"
)

// InMemoryEventBus es una implementación en memoria del EventBus
type InMemoryEventBus struct {
	handlers map[string][]domain.EventHandler
	mutex    sync.RWMutex
}

// NewInMemoryEventBus crea una nueva instancia del EventBus en memoria
func NewInMemoryEventBus() domain.EventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]domain.EventHandler),
	}
}

// Publish publica un evento del dominio a todos los handlers suscritos
func (b *InMemoryEventBus) Publish(ctx context.Context, event domain.DomainEvent) error {
	b.mutex.RLock()
	handlers := b.handlers[event.EventType()]
	b.mutex.RUnlock()

	// Ejecutar handlers de forma asíncrona para no bloquear
	for _, handler := range handlers {
		go func(h domain.EventHandler) {
			if err := h.Handle(ctx, event); err != nil {
				// Log error pero no fallar la publicación
				// En producción, podrías usar un logger estructurado
			}
		}(handler)
	}

	return nil
}

// Subscribe suscribe un handler a un tipo de evento específico
func (b *InMemoryEventBus) Subscribe(eventType string, handler domain.EventHandler) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}
