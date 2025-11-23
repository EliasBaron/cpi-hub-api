package domain

import "time"

// DomainEvent es la interfaz base para todos los eventos del dominio
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}

// CommentCreatedEvent se emite cuando se crea un comentario en un post
type CommentCreatedEvent struct {
	CommentID   int
	PostID      int
	CreatedBy   int
	PostOwnerID int
	Content     string
	IsReply     bool
	ParentID    *int
	occurredAt  time.Time
}

func (e CommentCreatedEvent) EventType() string {
	return "comment_created"
}

func (e CommentCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// NewCommentCreatedEvent crea un nuevo evento de comentario creado
func NewCommentCreatedEvent(commentID, postID, createdBy, postOwnerID int, content string, isReply bool, parentID *int, occurredAt time.Time) CommentCreatedEvent {
	return CommentCreatedEvent{
		CommentID:   commentID,
		PostID:      postID,
		CreatedBy:   createdBy,
		PostOwnerID: postOwnerID,
		Content:     content,
		IsReply:     isReply,
		ParentID:    parentID,
		occurredAt:  occurredAt,
	}
}

// CommentReplyCreatedEvent se emite cuando se crea una respuesta a un comentario
type CommentReplyCreatedEvent struct {
	CommentID     int
	PostID        int
	CreatedBy     int
	ParentOwnerID int
	ParentID      int
	Content       string
	occurredAt    time.Time
}

func (e CommentReplyCreatedEvent) EventType() string {
	return "comment_reply_created"
}

func (e CommentReplyCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// NewCommentReplyCreatedEvent crea un nuevo evento de respuesta a comentario creada
func NewCommentReplyCreatedEvent(commentID, postID, createdBy, parentOwnerID, parentID int, content string, occurredAt time.Time) CommentReplyCreatedEvent {
	return CommentReplyCreatedEvent{
		CommentID:     commentID,
		PostID:        postID,
		CreatedBy:     createdBy,
		ParentOwnerID: parentOwnerID,
		ParentID:      parentID,
		Content:       content,
		occurredAt:    occurredAt,
	}
}

// ReactionCreatedEvent se emite cuando se crea una reacción
type ReactionCreatedEvent struct {
	ReactionID  string
	UserID      int
	OwnerUserID int
	EntityType  string
	EntityID    int
	Action      string
	PostID      *int
	occurredAt  time.Time
}

func (e ReactionCreatedEvent) EventType() string {
	return "reaction_created"
}

func (e ReactionCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// NewReactionCreatedEvent crea un nuevo evento de reacción creada
func NewReactionCreatedEvent(reactionID string, userID, ownerUserID int, entityType string, entityID int, action string, postID *int, occurredAt time.Time) ReactionCreatedEvent {
	return ReactionCreatedEvent{
		ReactionID:  reactionID,
		UserID:      userID,
		OwnerUserID: ownerUserID,
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      action,
		PostID:      postID,
		occurredAt:  occurredAt,
	}
}
