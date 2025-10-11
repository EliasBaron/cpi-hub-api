package dto

import "net/http"

type EventsConnectionParams struct {
	UserID  int
	SpaceID int
	Writer  http.ResponseWriter
	Request *http.Request
}

type EventsBroadcastParams struct {
	SpaceID  int    `json:"space_id" binding:"required"`
	UserID   int    `json:"user_id" binding:"required"`
	Message  string `json:"message" binding:"required"`
	Username string `json:"username" binding:"required"`
}
