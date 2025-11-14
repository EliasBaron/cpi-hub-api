package events

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/dto"
	"cpi-hub-api/pkg/apperror"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// NotificationManager implementa la interfaz NotificationManager del dominio
type NotificationManager struct {
	connections map[int]*websocket.Conn // user_id -> connection
	mutex       sync.RWMutex
	upgrader    websocket.Upgrader
	config      *WebSocketConfig
}

func NewNotificationManager() domain.NotificationManager {
	return &NotificationManager{
		connections: make(map[int]*websocket.Conn),
		config:      DefaultWebSocketConfig(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  int(DefaultWebSocketConfig().MaxMessageSize),
			WriteBufferSize: int(DefaultWebSocketConfig().MaxMessageSize),
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (nm *NotificationManager) HandleConnection(params domain.HandleNotificationConnectionParams) error {
	if len(nm.connections) >= nm.config.MaxConnections {
		return apperror.NewInternalServer("maximum connections reached", nil, "notification_manager.go:HandleConnection")
	}

	conn, err := nm.upgrader.Upgrade(params.Writer, params.Request, nil)
	if err != nil {
		return apperror.NewInternalServer("error upgrading connection", err, "notification_manager.go:HandleConnection")
	}

	nm.mutex.Lock()
	if existingConn, exists := nm.connections[params.UserID]; exists {
		existingConn.Close()
	}
	nm.connections[params.UserID] = conn
	nm.mutex.Unlock()

	go nm.handleMessages(params.UserID, conn)

	return nil
}

// handleMessages maneja los mensajes entrantes de una conexión
func (nm *NotificationManager) handleMessages(userID int, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		nm.removeConnection(userID)
	}()

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(nm.config.PongWait))
		return nil
	})

	ticker := time.NewTicker(nm.config.GetPingPeriod())
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			nm.mutex.RLock()
			conn, exists := nm.connections[userID]
			nm.mutex.RUnlock()

			if !exists {
				return
			}

			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

	conn.SetReadDeadline(time.Now().Add(nm.config.PongWait))

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			break
		}
		conn.SetReadDeadline(time.Now().Add(nm.config.PongWait))
	}
}

func (nm *NotificationManager) removeConnection(userID int) {
	nm.mutex.Lock()
	delete(nm.connections, userID)
	nm.mutex.Unlock()
}

// BroadcastToUser envía una notificación a un usuario específico si está conectado
func (nm *NotificationManager) BroadcastToUser(userID int, notification *domain.Notification) error {
	nm.mutex.RLock()
	conn, exists := nm.connections[userID]
	nm.mutex.RUnlock()

	if !exists {
		return nil
	}

	notificationMessage := dto.ToNotificationMessageDTO(notification)

	messageBytes, err := json.Marshal(notificationMessage)
	if err != nil {
		return apperror.NewInternalServer("error marshaling notification", err, "notification_manager.go:BroadcastToUser")
	}

	conn.SetWriteDeadline(time.Now().Add(nm.config.WriteWait))
	if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		nm.removeConnection(userID)
		return apperror.NewInternalServer("error writing notification", err, "notification_manager.go:BroadcastToUser")
	}

	return nil
}
