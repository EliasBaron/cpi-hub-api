package handlers

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	eventsUsecase "cpi-hub-api/internal/core/usecase/events"
	"log"
)

// CommentEventHandler maneja eventos relacionados con comentarios
type CommentEventHandler struct {
	eventEmitter eventsUsecase.EventEmitter
}

// NewCommentEventHandler crea una nueva instancia del handler de comentarios
func NewCommentEventHandler(eventEmitter eventsUsecase.EventEmitter) domain.EventHandler {
	return &CommentEventHandler{
		eventEmitter: eventEmitter,
	}
}

// Handle procesa los eventos de comentarios y los convierte a eventos de WebSocket
func (h *CommentEventHandler) Handle(ctx context.Context, event domain.DomainEvent) error {
	switch e := event.(type) {
	case domain.CommentCreatedEvent:
		return h.handleCommentCreated(ctx, e)
	case domain.CommentReplyCreatedEvent:
		return h.handleCommentReplyCreated(ctx, e)
	default:
		// Ignorar eventos que no son de comentarios
		return nil
	}
}

func (h *CommentEventHandler) handleCommentCreated(ctx context.Context, event domain.CommentCreatedEvent) error {
	if h.eventEmitter == nil {
		return nil
	}

	wsEvent := &domain.Event{
		Type:         "comment_created",
		UserID:       event.CreatedBy,
		TargetUserID: event.PostOwnerID,
		Metadata: map[string]interface{}{
			"post_id":         event.PostID,
			"comment_id":      event.CommentID,
			"comment_content": event.Content,
			"is_reply":        event.IsReply,
		},
		Timestamp: event.OccurredAt(),
	}

	err := h.eventEmitter.EmitEvent(ctx, wsEvent)
	if err != nil {
		log.Printf("Error emitting comment created event: %v", err)
		return err
	}

	return nil
}

func (h *CommentEventHandler) handleCommentReplyCreated(ctx context.Context, event domain.CommentReplyCreatedEvent) error {
	if h.eventEmitter == nil {
		return nil
	}

	wsEvent := &domain.Event{
		Type:         "comment_reply_created",
		UserID:       event.CreatedBy,
		TargetUserID: event.ParentOwnerID,
		Metadata: map[string]interface{}{
			"post_id":           event.PostID,
			"comment_id":        event.CommentID,
			"parent_comment_id": event.ParentID,
			"comment_content":   event.Content,
			"is_reply":          true,
		},
		Timestamp: event.OccurredAt(),
	}

	err := h.eventEmitter.EmitEvent(ctx, wsEvent)
	if err != nil {
		log.Printf("Error emitting comment reply created event: %v", err)
		return err
	}

	return nil
}
