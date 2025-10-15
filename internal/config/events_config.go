package config

import "time"

// EventsConfig configuración para el sistema de eventos en tiempo real
type EventsConfig struct {
	// Configuración de WebSocket
	WebSocket WebSocketConfig `json:"websocket"`

	// Configuración de mensajes
	Messages MessagesConfig `json:"messages"`

	// Configuración de conexiones
	Connections ConnectionsConfig `json:"connections"`

	// Configuración de base de datos
	Database DatabaseConfig `json:"database"`
}

// WebSocketConfig configuración específica de WebSocket
type WebSocketConfig struct {
	// Tiempo de espera para escribir un mensaje al cliente
	WriteWait time.Duration `json:"write_wait" default:"10s"`

	// Tiempo de espera para leer el siguiente pong del cliente
	PongWait time.Duration `json:"pong_wait" default:"60s"`

	// Enviar pings al cliente con este período
	PingPeriod time.Duration `json:"ping_period" default:"54s"`

	// Tamaño máximo de mensaje permitido del cliente
	MaxMessageSize int64 `json:"max_message_size" default:"512"`

	// Tamaño del buffer de lectura
	ReadBufferSize int `json:"read_buffer_size" default:"1024"`

	// Tamaño del buffer de escritura
	WriteBufferSize int `json:"write_buffer_size" default:"1024"`

	// Habilitar compresión
	EnableCompression bool `json:"enable_compression" default:"false"`

	// Validar origen (en producción debe ser true)
	CheckOrigin bool `json:"check_origin" default:"false"`
}

// MessagesConfig configuración de mensajes
type MessagesConfig struct {
	// Longitud máxima del contenido del mensaje
	MaxContentLength int `json:"max_content_length" default:"1000"`

	// Longitud mínima del contenido del mensaje
	MinContentLength int `json:"min_content_length" default:"1"`

	// Límite de mensajes por espacio para historial
	HistoryLimit int `json:"history_limit" default:"50"`

	// Tiempo de retención de mensajes (0 = indefinido)
	RetentionPeriod time.Duration `json:"retention_period" default:"0"`

	// Habilitar persistencia de mensajes
	EnablePersistence bool `json:"enable_persistence" default:"true"`
}

// ConnectionsConfig configuración de conexiones
type ConnectionsConfig struct {
	// Tamaño del canal de envío por cliente
	SendChannelSize int `json:"send_channel_size" default:"256"`

	// Tamaño del canal de registro del hub
	RegisterChannelSize int `json:"register_channel_size" default:"256"`

	// Tamaño del canal de desregistro del hub
	UnregisterChannelSize int `json:"unregister_channel_size" default:"256"`

	// Tamaño del canal de broadcast del hub
	BroadcastChannelSize int `json:"broadcast_channel_size" default:"256"`

	// Tamaño del canal de broadcast por espacio
	SpaceBroadcastChannelSize int `json:"space_broadcast_channel_size" default:"256"`

	// Tiempo de timeout para conexiones inactivas
	InactiveTimeout time.Duration `json:"inactive_timeout" default:"300s"`

	// Intervalo de limpieza de conexiones inactivas
	CleanupInterval time.Duration `json:"cleanup_interval" default:"60s"`
}

// DatabaseConfig configuración de base de datos para eventos
type DatabaseConfig struct {
	// Habilitar logging de eventos del sistema
	EnableSystemEvents bool `json:"enable_system_events" default:"true"`

	// Habilitar tracking de conexiones activas
	EnableActiveConnections bool `json:"enable_active_connections" default:"true"`

	// Intervalo de limpieza de conexiones desconectadas
	CleanupDisconnectedInterval time.Duration `json:"cleanup_disconnected_interval" default:"300s"`

	// Tiempo después del cual considerar una conexión como desconectada
	DisconnectedThreshold time.Duration `json:"disconnected_threshold" default:"120s"`
}

// DefaultEventsConfig retorna la configuración por defecto
func DefaultEventsConfig() *EventsConfig {
	return &EventsConfig{
		WebSocket: WebSocketConfig{
			WriteWait:         10 * time.Second,
			PongWait:          60 * time.Second,
			PingPeriod:        54 * time.Second,
			MaxMessageSize:    512,
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			EnableCompression: false,
			CheckOrigin:       false, // Solo para desarrollo
		},
		Messages: MessagesConfig{
			MaxContentLength:  1000,
			MinContentLength:  1,
			HistoryLimit:      50,
			RetentionPeriod:   0, // Sin límite de tiempo
			EnablePersistence: true,
		},
		Connections: ConnectionsConfig{
			SendChannelSize:           256,
			RegisterChannelSize:       256,
			UnregisterChannelSize:     256,
			BroadcastChannelSize:      256,
			SpaceBroadcastChannelSize: 256,
			InactiveTimeout:           300 * time.Second,
			CleanupInterval:           60 * time.Second,
		},
		Database: DatabaseConfig{
			EnableSystemEvents:          true,
			EnableActiveConnections:     true,
			CleanupDisconnectedInterval: 300 * time.Second,
			DisconnectedThreshold:       120 * time.Second,
		},
	}
}

// ProductionEventsConfig retorna la configuración para producción
func ProductionEventsConfig() *EventsConfig {
	config := DefaultEventsConfig()

	// Configuración específica para producción
	config.WebSocket.CheckOrigin = true
	config.WebSocket.EnableCompression = true
	config.WebSocket.MaxMessageSize = 1024

	config.Messages.MaxContentLength = 2000
	config.Messages.HistoryLimit = 100

	config.Connections.SendChannelSize = 512
	config.Connections.InactiveTimeout = 600 * time.Second

	return config
}

// DevelopmentEventsConfig retorna la configuración para desarrollo
func DevelopmentEventsConfig() *EventsConfig {
	config := DefaultEventsConfig()

	// Configuración específica para desarrollo
	config.WebSocket.CheckOrigin = false
	config.WebSocket.EnableCompression = false

	config.Messages.MaxContentLength = 500
	config.Messages.HistoryLimit = 20

	config.Connections.InactiveTimeout = 120 * time.Second

	return config
}
