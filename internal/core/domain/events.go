package domain

import (
	"errors"
	"io"
	"net/http"
	"time"
)

// MessageType representa el tipo de mensaje de eventos en tiempo real
type MessageType string

const (
	MessageTypeChat  MessageType = "chat"
	MessageTypeJoin  MessageType = "join"
	MessageTypeLeave MessageType = "leave"
	MessageTypeError MessageType = "error"
	MessageTypePing  MessageType = "ping"
	MessageTypePong  MessageType = "pong"
)

// EventMessage representa un mensaje de eventos en tiempo real genérico
type EventMessage struct {
	Type      MessageType `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	UserID    int         `json:"user_id,omitempty"`
	SpaceID   int         `json:"space_id,omitempty"`
	Username  string      `json:"username,omitempty"`
	Image     string      `json:"image,omitempty"`
}

// ChatMessage representa un mensaje de chat específico
type ChatMessage struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	SpaceID   int       `json:"space_id"`
	Timestamp time.Time `json:"timestamp"`
	Image     string    `json:"image"`
}

// JoinMessage representa un mensaje de unión a un espacio
type JoinMessage struct {
	SpaceID  int    `json:"space_id"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

// LeaveMessage representa un mensaje de salida de un espacio
type LeaveMessage struct {
	SpaceID  int    `json:"space_id"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

// ErrorMessage representa un mensaje de error
type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Client representa un cliente conectado a eventos en tiempo real
type Client struct {
	ID       string
	UserID   int
	SpaceID  int
	Send     chan []byte
	Hub      *Hub
	Conn     EventConnection
	Username string
}

// Hub mantiene el conjunto de clientes activos y los mensajes de difusión
type Hub struct {
	// Clientes registrados
	Clients map[*Client]bool

	// Canal para registrar clientes
	Register chan *Client

	// Canal para desregistrar clientes
	Unregister chan *Client

	// Canal para difundir mensajes a todos los clientes
	Broadcast chan []byte

	// Canal para enviar mensajes a un espacio específico
	SpaceBroadcast chan SpaceMessage
}

// SpaceMessage representa un mensaje dirigido a un espacio específico
type SpaceMessage struct {
	SpaceID int
	Message []byte
}

// EventConnection define la interfaz para conexiones de eventos en tiempo real
// Esta interfaz debe ser implementada por *websocket.Conn de Gorilla WebSocket
type EventConnection interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
	SetReadLimit(limit int64)
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	SetPongHandler(handler func(string) error)
	NextWriter(messageType int) (io.WriteCloser, error)
}

// Errores específicos de eventos en tiempo real
var (
	ErrEmptyMessage   = errors.New("el mensaje no puede estar vacío")
	ErrMessageTooLong = errors.New("el mensaje es demasiado largo")
	ErrUnauthorized   = errors.New("no autorizado para realizar esta acción")
)

type UserStatus string

const (
	UserStatusOnline  UserStatus = "online"
	UserStatusOffline UserStatus = "offline"
)

type UserConnectionMessage struct {
	Type      string     `json:"type"`
	UserID    int        `json:"user_id"`
	Status    UserStatus `json:"status"`
	Username  string     `json:"username"`
	Timestamp string     `json:"timestamp"`
}

type HandleUserConnectionParams struct {
	UserID  int
	Writer  http.ResponseWriter
	Request *http.Request
}

type UserConnectionManager interface {
	HandleConnection(params HandleUserConnectionParams) error
}

type HandleNotificationConnectionParams struct {
	UserID  int
	Writer  http.ResponseWriter
	Request *http.Request
}

type NotificationManager interface {
	HandleConnection(params HandleNotificationConnectionParams) error
	BroadcastEvent(userID int, event *Event) error
}
