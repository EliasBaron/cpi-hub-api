package events

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/dto"
	websocketAdapter "cpi-hub-api/internal/infrastructure/adapters/websocket"
	"cpi-hub-api/pkg/apperror"
	"cpi-hub-api/pkg/helpers"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// EventsUsecase maneja la lógica de negocio para eventos en tiempo real
type EventsUsecase struct {
	hubManager      *HubManager
	userConnManager domain.UserConnectionManager
	repository      domain.EventsRepository
	userRepository  domain.UserRepository
	spaceRepository domain.SpaceRepository
	config          *WebSocketConfig
}

// NewEventsUsecase crea una nueva instancia del EventsUsecase
func NewEventsUsecase(
	hubManager *HubManager,
	userConnManager domain.UserConnectionManager,
	repository domain.EventsRepository,
	userRepository domain.UserRepository,
	spaceRepository domain.SpaceRepository,
) *EventsUsecase {
	return &EventsUsecase{
		hubManager:      hubManager,
		userConnManager: userConnManager,
		repository:      repository,
		userRepository:  userRepository,
		spaceRepository: spaceRepository,
		config:          DefaultWebSocketConfig(),
	}
}

// HandleConnection maneja toda la lógica de conexión WebSocket
func (u *EventsUsecase) HandleConnection(params dto.EventsConnectionParams) error {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  int(u.config.MaxMessageSize),
		WriteBufferSize: int(u.config.MaxMessageSize),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(params.Writer, params.Request, nil)
	if err != nil {
		return apperror.NewInternalServer("Error upgrading WebSocket connection", err, "events_usecase.go:HandleConnection")
	}

	// Crear wrapper de WebSocket
	wsConn := websocketAdapter.NewWebSocketWrapper(conn)

	// Crear cliente
	client := u.CreateClient(params.UserID, params.SpaceID, params.Username, wsConn)

	// Registrar cliente en el hub
	u.RegisterClient(client)

	// Crear y iniciar client manager
	clientManager := NewClientManager(client)
	go clientManager.WritePump()
	go clientManager.ReadPump()

	return nil
}

// CreateClient crea un cliente para el hub (usado por la capa de infraestructura)
func (u *EventsUsecase) CreateClient(userID, spaceID int, username string, conn domain.EventConnection) *domain.Client {
	return &domain.Client{
		ID:       u.generateClientID(userID, spaceID),
		UserID:   userID,
		SpaceID:  spaceID,
		Username: username,
		Send:     make(chan []byte, u.config.SendBufferSize),
		Hub:      u.hubManager.GetHub(),
		Conn:     conn,
	}
}

// RegisterClient registra un cliente en el hub
func (u *EventsUsecase) RegisterClient(client *domain.Client) {
	u.hubManager.GetHub().Register <- client
}

// Broadcast envía un mensaje a todos los usuarios de un espacio específico
func (u *EventsUsecase) Broadcast(dto dto.EventsBroadcastParams) (*domain.ChatMessage, error) {
	return u.broadcastMessage(dto)
}

// BroadcastToSpace es un alias para mantener compatibilidad
func (u *EventsUsecase) BroadcastToSpace(dto dto.EventsBroadcastParams) (*domain.ChatMessage, error) {
	return u.broadcastMessage(dto)
}

// broadcastMessage maneja la lógica común de difusión de mensajes
func (u *EventsUsecase) broadcastMessage(dto dto.EventsBroadcastParams) (*domain.ChatMessage, error) {
	if err := u.validateMessageContent(dto.Message); err != nil {
		return nil, err
	}

	chatMsg := &domain.ChatMessage{
		ID:        helpers.NewULID(),
		Content:   dto.Message,
		UserID:    dto.UserID,
		Username:  dto.Username,
		SpaceID:   dto.SpaceID,
		Image:     dto.Image,
		Timestamp: helpers.GetTime(),
	}

	if err := u.repository.SaveMessage(chatMsg); err != nil {
		return nil, err
	}

	u.hubManager.BroadcastChatMessage(chatMsg)
	return chatMsg, nil
}

func (u *EventsUsecase) validateMessageContent(content string) error {
	if content == "" {
		return domain.ErrEmptyMessage
	}

	if len(content) > 1000 {
		return domain.ErrMessageTooLong
	}

	return nil
}

func (u *EventsUsecase) generateClientID(userID, spaceID int) string {
	return fmt.Sprintf("%d-%d-%s", userID, spaceID, helpers.NewULID())
}

// HandleUserConnection maneja la conexión de usuario para el estado online/offline
func (u *EventsUsecase) HandleUserConnection(params dto.HandleUserConnectionParams) error {
	handleUserConnectionParams := domain.HandleUserConnectionParams{
		UserID:  params.UserID,
		Writer:  params.Writer,
		Request: params.Request,
	}

	return u.userConnManager.HandleConnection(handleUserConnectionParams)
}
