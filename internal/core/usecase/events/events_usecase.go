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

type EventsUsecase struct {
	hubManager          *HubManager
	userConnManager     domain.UserConnectionManager
	notificationManager domain.NotificationManager
	repository          domain.EventsRepository
	userRepository      domain.UserRepository
	spaceRepository     domain.SpaceRepository
	config              *WebSocketConfig
}

func NewEventsUsecase(
	hubManager *HubManager,
	userConnManager domain.UserConnectionManager,
	notificationManager domain.NotificationManager,
	repository domain.EventsRepository,
	userRepository domain.UserRepository,
	spaceRepository domain.SpaceRepository,
) *EventsUsecase {
	return &EventsUsecase{
		hubManager:          hubManager,
		userConnManager:     userConnManager,
		notificationManager: notificationManager,
		repository:          repository,
		userRepository:      userRepository,
		spaceRepository:     spaceRepository,
		config:              DefaultWebSocketConfig(),
	}
}

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

	wsConn := websocketAdapter.NewWebSocketWrapper(conn)

	client := u.CreateClient(params.UserID, params.SpaceID, params.Username, wsConn)

	u.RegisterClient(client)

	clientManager := NewClientManager(client)
	go clientManager.WritePump()
	go clientManager.ReadPump()

	return nil
}

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

func (u *EventsUsecase) RegisterClient(client *domain.Client) {
	u.hubManager.GetHub().Register <- client
}

func (u *EventsUsecase) Broadcast(dto dto.EventsBroadcastParams) (*domain.ChatMessage, error) {
	return u.broadcastMessage(dto)
}

func (u *EventsUsecase) BroadcastToSpace(dto dto.EventsBroadcastParams) (*domain.ChatMessage, error) {
	return u.broadcastMessage(dto)
}

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

func (u *EventsUsecase) HandleUserConnection(params dto.HandleUserConnectionParams) error {
	handleUserConnectionParams := domain.HandleUserConnectionParams{
		UserID:  params.UserID,
		Writer:  params.Writer,
		Request: params.Request,
	}

	return u.userConnManager.HandleConnection(handleUserConnectionParams)
}

func (u *EventsUsecase) HandleNotificationConnection(params dto.HandleNotificationConnectionParams) error {
	handleNotificationConnectionParams := domain.HandleNotificationConnectionParams{
		UserID:  params.UserID,
		Writer:  params.Writer,
		Request: params.Request,
	}

	return u.notificationManager.HandleConnection(handleNotificationConnectionParams)
}
