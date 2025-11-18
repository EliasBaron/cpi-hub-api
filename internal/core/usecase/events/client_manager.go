package events

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/pkg/helpers"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// ClientManager maneja las operaciones del cliente
type ClientManager struct {
	client *domain.Client
	config *WebSocketConfig
}

// NewClientManager crea una nueva instancia del ClientManager
func NewClientManager(client *domain.Client) *ClientManager {
	return &ClientManager{
		client: client,
		config: DefaultWebSocketConfig(),
	}
}

// readPump bombea mensajes desde la conexión WebSocket al hub
func (cm *ClientManager) ReadPump() {
	defer func() {
		cm.client.Hub.Unregister <- cm.client
		cm.client.Conn.Close()
	}()

	cm.client.Conn.SetReadLimit(cm.config.MaxMessageSize)
	cm.client.Conn.SetReadDeadline(helpers.GetTime().Add(cm.config.PongWait))
	cm.client.Conn.SetPongHandler(func(string) error {
		cm.client.Conn.SetReadDeadline(helpers.GetTime().Add(cm.config.PongWait))
		return nil
	})

	for {
		_, messageBytes, err := cm.client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error de WebSocket: %v", err)
			}
			break
		}

		// Procesar el mensaje recibido
		cm.handleMessage(messageBytes)
	}
}

// writePump bombea mensajes desde el hub a la conexión WebSocket
func (cm *ClientManager) WritePump() {
	ticker := time.NewTicker(cm.config.GetPingPeriod())
	defer func() {
		ticker.Stop()
		cm.client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-cm.client.Send:
			cm.client.Conn.SetWriteDeadline(helpers.GetTime().Add(cm.config.WriteWait))
			if !ok {
				// El hub cerró el canal
				cm.client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := cm.client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Agregar mensajes en cola al mensaje actual
			n := len(cm.client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-cm.client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			cm.client.Conn.SetWriteDeadline(helpers.GetTime().Add(cm.config.WriteWait))
			if err := cm.client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage procesa los mensajes recibidos del cliente
func (cm *ClientManager) handleMessage(messageBytes []byte) {
	var wsMsg domain.EventMessage
	if err := json.Unmarshal(messageBytes, &wsMsg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		cm.sendError("invalid_message_format", "Invalid message format")
		return
	}

	// Establecer información del cliente en el mensaje
	wsMsg.UserID = cm.client.UserID
	wsMsg.SpaceID = cm.client.SpaceID
	wsMsg.Timestamp = helpers.GetTime()

	// Procesar según el tipo de mensaje
	switch wsMsg.Type {
	case domain.MessageTypeChat:
		cm.handleChatMessage(wsMsg)
	case domain.MessageTypePing:
		cm.handlePing(wsMsg)
	default:
		log.Printf("Unknown message type: %s", wsMsg.Type)
		cm.sendError("unknown_message_type", "Unknown message type")
	}
}

// handleChatMessage procesa mensajes de chat
func (cm *ClientManager) handleChatMessage(wsMsg domain.EventMessage) {
	// Convertir data a ChatMessage
	dataBytes, err := json.Marshal(wsMsg.Data)
	if err != nil {
		log.Printf("Error marshaling chat data: %v", err)
		cm.sendError("invalid_chat_data", "Invalid chat data")
		return
	}

	var chatMsg domain.ChatMessage
	if err := json.Unmarshal(dataBytes, &chatMsg); err != nil {
		log.Printf("Error unmarshaling chat message: %v", err)
		cm.sendError("invalid_chat_message", "Invalid chat message")
		return
	}

	// Validar el mensaje
	if err := cm.validateChatMessage(&chatMsg); err != nil {
		cm.sendError("validation_error", err.Error())
		return
	}

	// Establecer información del cliente
	chatMsg.UserID = cm.client.UserID
	chatMsg.SpaceID = cm.client.SpaceID
	chatMsg.Timestamp = helpers.GetTime()
	chatMsg.ID = helpers.NewULID()

	// Difundir el mensaje al espacio
	cm.broadcastChatMessage(&chatMsg)
}

// validateChatMessage valida un mensaje de chat
func (cm *ClientManager) validateChatMessage(chatMsg *domain.ChatMessage) error {
	if chatMsg.Content == "" {
		return domain.ErrEmptyMessage
	}
	if len(chatMsg.Content) > 1000 {
		return domain.ErrMessageTooLong
	}
	return nil
}

// handlePing procesa mensajes ping
func (cm *ClientManager) handlePing(wsMsg domain.EventMessage) {
	pongMsg := domain.EventMessage{
		Type:      domain.MessageTypePong,
		Data:      map[string]string{"message": "pong"},
		Timestamp: helpers.GetTime(),
		UserID:    cm.client.UserID,
		SpaceID:   cm.client.SpaceID,
	}

	cm.sendMessage(pongMsg)
}

// sendError envía un mensaje de error al cliente
func (cm *ClientManager) sendError(code, message string) {
	errorMsg := domain.EventMessage{
		Type:      domain.MessageTypeError,
		Data:      domain.ErrorMessage{Code: code, Message: message},
		Timestamp: helpers.GetTime(),
		UserID:    cm.client.UserID,
		SpaceID:   cm.client.SpaceID,
	}

	cm.sendMessage(errorMsg)
}

// broadcastChatMessage difunde un mensaje de chat al espacio
func (cm *ClientManager) broadcastChatMessage(chatMsg *domain.ChatMessage) {
	wsMsg := domain.EventMessage{
		Type:      domain.MessageTypeChat,
		Data:      chatMsg,
		Timestamp: helpers.GetTime(),
		UserID:    chatMsg.UserID,
		SpaceID:   chatMsg.SpaceID,
	}

	messageBytes, err := json.Marshal(wsMsg)
	if err != nil {
		log.Printf("Error marshaling chat message: %v", err)
		return
	}

	// Enviar mensaje a todos los clientes del mismo espacio
	spaceMsg := domain.SpaceMessage{
		SpaceID: chatMsg.SpaceID,
		Message: messageBytes,
	}

	select {
	case cm.client.Hub.SpaceBroadcast <- spaceMsg:
		log.Printf("Chat message broadcasted: user %d in space %d: %s",
			chatMsg.UserID, chatMsg.SpaceID, chatMsg.Content)
	default:
		log.Printf("Could not broadcast message to space %d", chatMsg.SpaceID)
	}
}

// sendMessage envía un mensaje al cliente
func (cm *ClientManager) sendMessage(wsMsg domain.EventMessage) {
	messageBytes, err := json.Marshal(wsMsg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case cm.client.Send <- messageBytes:
	default:
		close(cm.client.Send)
		delete(cm.client.Hub.Clients, cm.client)
	}
}
