package events

import "time"

// WebSocketConfig contiene todas las configuraciones para WebSocket
type WebSocketConfig struct {
	// Timeouts
	WriteWait  time.Duration
	PongWait   time.Duration
	PingPeriod time.Duration

	// Message limits
	MaxMessageSize int64
	SendBufferSize int

	// Connection limits
	MaxConnections int
}

// DefaultWebSocketConfig retorna la configuración por defecto para WebSocket
func DefaultWebSocketConfig() *WebSocketConfig {
	return &WebSocketConfig{
		WriteWait:      10 * time.Second,
		PongWait:       60 * time.Second,
		PingPeriod:     54 * time.Second, // (pongWait * 9) / 10
		MaxMessageSize: 512,
		SendBufferSize: 256,
		MaxConnections: 50,
	}
}

// GetPingPeriod calcula el período de ping basado en pongWait
func (c *WebSocketConfig) GetPingPeriod() time.Duration {
	return (c.PongWait * 9) / 10
}
