package events

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/pkg/helpers"
	"encoding/json"
	"log"
)

// HubManager maneja las operaciones del hub
type HubManager struct {
	hub *domain.Hub
}

// NewHubManager crea una nueva instancia del HubManager
func NewHubManager() *HubManager {
	return &HubManager{
		hub: &domain.Hub{
			Clients:        make(map[*domain.Client]bool),
			Register:       make(chan *domain.Client),
			Unregister:     make(chan *domain.Client),
			Broadcast:      make(chan []byte),
			SpaceBroadcast: make(chan domain.SpaceMessage, 100), // buffer de 100
		},
	}
}

// GetHub devuelve el hub subyacente
func (hm *HubManager) GetHub() *domain.Hub {
	return hm.hub
}

// Run inicia el hub y maneja los canales
func (hm *HubManager) Run() {
	for {
		select {
		case client := <-hm.hub.Register:
			hm.hub.Clients[client] = true
			log.Printf("Cliente %s (ID %d) conectado al espacio %d", client.Username, client.UserID, client.SpaceID)

			// Enviar mensaje de bienvenida
			welcomeMsg := domain.EventMessage{
				Type:      domain.MessageTypeJoin,
				Data:      domain.JoinMessage{SpaceID: client.SpaceID, UserID: client.UserID},
				Timestamp: helpers.GetTime(),
				UserID:    client.UserID,
				SpaceID:   client.SpaceID,
				Username:  client.Username,
			}
			hm.broadcastToSpace(client.SpaceID, welcomeMsg)

		case client := <-hm.hub.Unregister:
			if _, ok := hm.hub.Clients[client]; ok {
				delete(hm.hub.Clients, client)
				close(client.Send)
				log.Printf("Cliente %d desconectado del espacio %d", client.UserID, client.SpaceID)

				// Enviar mensaje de despedida
				leaveMsg := domain.EventMessage{
					Type:      domain.MessageTypeLeave,
					Data:      domain.LeaveMessage{SpaceID: client.SpaceID, UserID: client.UserID},
					Timestamp: helpers.GetTime(),
					UserID:    client.UserID,
					SpaceID:   client.SpaceID,
				}
				hm.broadcastToSpace(client.SpaceID, leaveMsg)
			}

		case message := <-hm.hub.Broadcast:
			for client := range hm.hub.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(hm.hub.Clients, client)
				}
			}

		case spaceMsg := <-hm.hub.SpaceBroadcast:
			for client := range hm.hub.Clients {
				if client.SpaceID == spaceMsg.SpaceID {
					select {
					case client.Send <- spaceMsg.Message:
					default:
						close(client.Send)
						delete(hm.hub.Clients, client)
					}
				}
			}
		}
	}
}

// broadcastToSpace envía un mensaje a todos los clientes de un espacio específico
func (hm *HubManager) broadcastToSpace(spaceID int, message domain.EventMessage) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	spaceMsg := domain.SpaceMessage{
		SpaceID: spaceID,
		Message: messageBytes,
	}

	select {
	case hm.hub.SpaceBroadcast <- spaceMsg:
	default:
		log.Printf("No se pudo enviar mensaje al espacio %s", spaceID)
	}
}

// BroadcastChatMessage envía un mensaje de chat a un espacio específico
func (hm *HubManager) BroadcastChatMessage(chatMsg *domain.ChatMessage) {
	wsMsg := domain.EventMessage{
		Type:      domain.MessageTypeChat,
		Data:      chatMsg,
		Timestamp: helpers.GetTime(),
		UserID:    chatMsg.UserID,
		SpaceID:   chatMsg.SpaceID,
	}

	hm.broadcastToSpace(chatMsg.SpaceID, wsMsg)
}
