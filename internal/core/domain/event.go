package domain

import "time"

// Event representa un evento genérico emitido por el backend
type Event struct {
	Type         string                 `json:"type"`
	UserID       int                    `json:"user_id"`
	TargetUserID int                    `json:"target_user_id"`
	Metadata     map[string]interface{} `json:"metadata"`
	Timestamp    time.Time              `json:"timestamp"`
}

