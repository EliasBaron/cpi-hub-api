package domain

import (
	"context"
	"cpi-hub-api/internal/core/domain/criteria"
)

//go:generate mockgen -package=mock -source=./repositories.go -destination=./mock/repositories_mock.go
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Find(ctx context.Context, criteria *criteria.Criteria) (*User, error)
	Search(ctx context.Context, criteria *criteria.Criteria) ([]*User, error)
	Count(ctx context.Context, criteria *criteria.Criteria) (int, error)
	Update(ctx context.Context, user *User) error
}

type SpaceRepository interface {
	Create(ctx context.Context, space *Space) error
	Find(ctx context.Context, criteria *criteria.Criteria) (*Space, error)
	FindByIDs(ctx context.Context, ids []int) ([]*Space, error)
	FindAll(ctx context.Context, criteria *criteria.Criteria) ([]*Space, error)
	Update(ctx context.Context, space *Space) error
	Count(ctx context.Context, criteria *criteria.Criteria) (int, error)
}

type UserSpaceRepository interface {
	Update(ctx context.Context, userId int, spaceIDs []int, action string) error
	FindSpacesIDsByUserID(ctx context.Context, userID int) ([]int, error)
	FindUserIDsBySpaceID(ctx context.Context, spaceID int) ([]int, error)
	Exists(ctx context.Context, userId int, spaceId int) (bool, error)
	Count(ctx context.Context, criteria *criteria.Criteria) (int, error)
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	Find(ctx context.Context, criteria *criteria.Criteria) (*Post, error)
	Update(ctx context.Context, post *Post) error
	Search(ctx context.Context, criteria *criteria.Criteria) ([]*Post, error)
	Count(ctx context.Context, criteria *criteria.Criteria) (int, error)
	Delete(ctx context.Context, postID int) error
}

type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	Find(ctx context.Context, criteria *criteria.Criteria) (*CommentWithInfo, error)
	FindAll(ctx context.Context, criteria *criteria.Criteria) ([]*CommentWithInfo, error)
	Count(ctx context.Context, criteria *criteria.Criteria) (int, error)
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, commentID int) error
}

type EventsRepository interface {
	SaveMessage(message *ChatMessage) error
}

type SearchMessagesFilter struct {
	SpaceID       int
	Page          int
	PageSize      int
	OrderBy       string
	SortDirection string
}

type MessageRepository interface {
	SearchMessages(ctx context.Context, filters SearchMessagesFilter) ([]*ChatMessage, int, error)
}

type ReactionRepository interface {
	AddReaction(ctx context.Context, reaction *Reaction) error
	FindReaction(ctx context.Context, criteria *criteria.Criteria) (*Reaction, error)
	DeleteReaction(ctx context.Context, reactionID string) error
	UpdateReaction(ctx context.Context, reaction *Reaction) error
	CountReactions(ctx context.Context, criteria *criteria.Criteria) (int, error)
	// GetReactions(ctx context.Context, criteria *criteria.Criteria) ([]*Reaction, error)
}

type NotificationRepository interface {
	SaveNotification(ctx context.Context, notification *Notification) error
	GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]*Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
	MarkAllAsRead(ctx context.Context, userID int) error
	GetUnreadCount(ctx context.Context, userID int) (int, error)
}
