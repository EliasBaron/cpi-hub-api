package websocket

import (
	"cpi-hub-api/internal/core/domain"
	"io"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketWrapper implementa la interfaz EventConnection usando *websocket.Conn
type WebSocketWrapper struct {
	conn *websocket.Conn
}

// NewWebSocketWrapper crea un nuevo wrapper para *websocket.Conn
func NewWebSocketWrapper(conn *websocket.Conn) domain.EventConnection {
	return &WebSocketWrapper{
		conn: conn,
	}
}

// ReadMessage implementa el método ReadMessage de la interfaz
func (w *WebSocketWrapper) ReadMessage() (messageType int, p []byte, err error) {
	return w.conn.ReadMessage()
}

// WriteMessage implementa el método WriteMessage de la interfaz
func (w *WebSocketWrapper) WriteMessage(messageType int, data []byte) error {
	return w.conn.WriteMessage(messageType, data)
}

// Close implementa el método Close de la interfaz
func (w *WebSocketWrapper) Close() error {
	return w.conn.Close()
}

// SetReadLimit implementa el método SetReadLimit de la interfaz
func (w *WebSocketWrapper) SetReadLimit(limit int64) {
	w.conn.SetReadLimit(limit)
}

// SetReadDeadline implementa el método SetReadDeadline de la interfaz
func (w *WebSocketWrapper) SetReadDeadline(t time.Time) error {
	return w.conn.SetReadDeadline(t)
}

// SetWriteDeadline implementa el método SetWriteDeadline de la interfaz
func (w *WebSocketWrapper) SetWriteDeadline(t time.Time) error {
	return w.conn.SetWriteDeadline(t)
}

// SetPongHandler implementa el método SetPongHandler de la interfaz
func (w *WebSocketWrapper) SetPongHandler(handler func(string) error) {
	w.conn.SetPongHandler(handler)
}

// NextWriter implementa el método NextWriter de la interfaz
func (w *WebSocketWrapper) NextWriter(messageType int) (io.WriteCloser, error) {
	return w.conn.NextWriter(messageType)
}
