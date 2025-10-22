package events

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/pkg/apperror"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// UserConnectionManager implementa la interfaz UserConnectionManager del dominio
type UserConnectionManager struct {
	connections map[int]*websocket.Conn // user_id -> connection
	userStatus  map[int]bool            // user_id -> online/offline
	mutex       sync.RWMutex
	upgrader    websocket.Upgrader
	config      *WebSocketConfig
}

func NewUserConnectionManager() domain.UserConnectionManager {
	return &UserConnectionManager{
		connections: make(map[int]*websocket.Conn),
		userStatus:  make(map[int]bool),
		config:      DefaultWebSocketConfig(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// HandleConnection maneja una nueva conexión de usuario
func (ucm *UserConnectionManager) HandleConnection(params domain.HandleUserConnectionParams) error {
	if len(ucm.connections) >= ucm.config.MaxConnections {
		return apperror.NewInternalServer("maximum connections reached", nil, "user_connection_manager.go:HandleConnection")
	}

	conn, err := ucm.upgrader.Upgrade(params.Writer, params.Request, nil)
	if err != nil {
		return apperror.NewInternalServer("error upgrading connection", err, "user_connection_manager.go:HandleConnection")
	}

	ucm.mutex.Lock()
	// Si el usuario ya está conectado, cerrar la conexión anterior
	if existingConn, exists := ucm.connections[params.UserID]; exists {
		existingConn.Close()
	}
	ucm.connections[params.UserID] = conn
	ucm.userStatus[params.UserID] = true
	ucm.mutex.Unlock()

	// Enviar mensaje inicial al usuario que se conecta
	ucm.sendInitialStatusMessage(params.UserID, conn)

	// Enviar lista de usuarios ya conectados al nuevo usuario
	ucm.sendConnectedUsersList(params.UserID, conn)

	// Notificar a otros usuarios que este usuario está online
	ucm.broadcastUserStatusToOthers(params.UserID, domain.UserStatusOnline)

	go ucm.handleMessages(params.UserID, conn)

	return nil
}

// handleMessages maneja los mensajes entrantes de una conexión
func (ucm *UserConnectionManager) handleMessages(userID int, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		ucm.removeConnection(userID)
	}()

	// Configurar ping/pong para detectar conexiones muertas
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(ucm.config.PongWait))
		return nil
	})

	// Enviar ping periódico para mantener la conexión viva
	ticker := time.NewTicker(ucm.config.GetPingPeriod())
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

	// Configurar timeout de lectura
	conn.SetReadDeadline(time.Now().Add(ucm.config.PongWait))

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Detectar diferentes tipos de errores de conexión
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Conexión cerrada inesperadamente
			}
			break
		}
		// Resetear el deadline después de cada mensaje
		conn.SetReadDeadline(time.Now().Add(ucm.config.PongWait))
	}
}

// removeConnection elimina una conexión de usuario
func (ucm *UserConnectionManager) removeConnection(userID int) {
	ucm.mutex.Lock()
	delete(ucm.connections, userID)
	ucm.userStatus[userID] = false
	ucm.mutex.Unlock()

	ucm.broadcastUserStatusToOthers(userID, domain.UserStatusOffline)
}

// sendInitialStatusMessage envía un mensaje inicial al usuario que se conecta
func (ucm *UserConnectionManager) sendInitialStatusMessage(userID int, conn *websocket.Conn) {
	message := domain.UserConnectionMessage{
		Type:      "user_status",
		UserID:    userID,
		Status:    domain.UserStatusOnline,
		Username:  ucm.getUsername(userID),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	err := conn.WriteJSON(message)
	if err != nil {
		// Si no se puede enviar el mensaje inicial, cerrar la conexión
		conn.Close()
		ucm.mutex.Lock()
		delete(ucm.connections, userID)
		ucm.mutex.Unlock()
	}
}

// broadcastUserStatusToOthers difunde el estado de un usuario a todos los demás conectados
func (ucm *UserConnectionManager) broadcastUserStatusToOthers(userID int, status domain.UserStatus) {
	message := domain.UserConnectionMessage{
		Type:      "user_status",
		UserID:    userID,
		Status:    status,
		Username:  ucm.getUsername(userID),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Crear una copia de las conexiones para evitar bloqueos largos
	ucm.mutex.RLock()
	connectionsToNotify := make(map[int]*websocket.Conn, len(ucm.connections))
	for id, conn := range ucm.connections {
		if id != userID {
			connectionsToNotify[id] = conn
		}
	}
	ucm.mutex.RUnlock()

	// Enviar mensajes de forma asíncrona para evitar bloqueos
	go func() {
		for id, conn := range connectionsToNotify {
			err := conn.WriteJSON(message)
			if err != nil {
				// Si hay error, remover la conexión
				ucm.mutex.Lock()
				delete(ucm.connections, id)
				ucm.mutex.Unlock()
			}
		}
	}()
}

// sendConnectedUsersList envía la lista de usuarios ya conectados al nuevo usuario
func (ucm *UserConnectionManager) sendConnectedUsersList(newUserID int, conn *websocket.Conn) {
	ucm.mutex.RLock()
	connectedUsers := make([]domain.UserConnectionMessage, 0)
	for userID := range ucm.connections {
		if userID != newUserID { // No incluir al usuario que se acaba de conectar
			connectedUsers = append(connectedUsers, domain.UserConnectionMessage{
				Type:      "user_status",
				UserID:    userID,
				Status:    domain.UserStatusOnline,
				Username:  ucm.getUsername(userID),
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
	}
	ucm.mutex.RUnlock()

	// Enviar cada usuario conectado como un mensaje separado
	for _, userMessage := range connectedUsers {
		err := conn.WriteJSON(userMessage)
		if err != nil {
			// Si no se puede enviar, cerrar la conexión
			conn.Close()
			ucm.mutex.Lock()
			delete(ucm.connections, newUserID)
			ucm.mutex.Unlock()
			break
		}
	}
}

// getUsername obtiene el nombre de usuario (placeholder - implementar según lógica de negocio)
func (ucm *UserConnectionManager) getUsername(userID int) string {
	// TODO: Implementar lógica para obtener el username del userID
	return "User" + strconv.Itoa(userID)
}
