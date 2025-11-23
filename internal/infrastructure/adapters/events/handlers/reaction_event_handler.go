package handlers

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	eventsUsecase "cpi-hub-api/internal/core/usecase/events"
	"log"
)

// ReactionEventHandler maneja eventos relacionados con reacciones
type ReactionEventHandler struct {
	eventEmitter eventsUsecase.EventEmitter
}

// NewReactionEventHandler crea una nueva instancia del handler de reacciones
func NewReactionEventHandler(eventEmitter eventsUsecase.EventEmitter) domain.EventHandler {
	return &ReactionEventHandler{
		eventEmitter: eventEmitter,
	}
}

// Handle procesa los eventos de reacciones y los convierte a eventos de WebSocket
func (h *ReactionEventHandler) Handle(ctx context.Context, event domain.DomainEvent) error {
	reactionEvent, ok := event.(domain.ReactionCreatedEvent)
	if !ok {
		// Ignorar eventos que no son de reacciones
		return nil
	}

	if h.eventEmitter == nil {
		return nil
	}

	wsEvent := &domain.Event{
		Type:         "reaction_created",
		UserID:       reactionEvent.UserID,
		TargetUserID: reactionEvent.OwnerUserID,
		Metadata: map[string]interface{}{
			"entity_type": reactionEvent.EntityType,
			"entity_id":   reactionEvent.EntityID,
			"action":      reactionEvent.Action,
		},
		Timestamp: reactionEvent.OccurredAt(),
	}

	// Agregar post_id si está disponible
	if reactionEvent.PostID != nil {
		wsEvent.Metadata["post_id"] = *reactionEvent.PostID
	}

	err := h.eventEmitter.EmitEvent(ctx, wsEvent)
	if err != nil {
		log.Printf("Error emitting reaction created event: %v", err)
		return err
	}

	return nil
}
