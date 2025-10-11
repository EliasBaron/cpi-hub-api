package dto

import "net/http"

type EventsConnectionParams struct {
	UserID  string
	SpaceID string
	Writer  http.ResponseWriter
	Request *http.Request
}

type EventsBroadcastParams struct {
	UserID   string `json:"user_id" binding:"required"`
	Message  string `json:"message" binding:"required"`
	Username string `json:"username" binding:"required"`
}
