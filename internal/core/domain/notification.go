package domain

import "time"

// Notification representa una notificación genérica
type Notification struct {
	ID          string
	Title       string
	Description string
	URL         *string
	To          int // user_id del destinatario
	Read        bool
	CreatedAt   time.Time
}
