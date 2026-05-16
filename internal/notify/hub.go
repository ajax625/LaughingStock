package notify

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type wsConn struct {
	send chan []byte
	conn *websocket.Conn
}

// Hub tracks active WebSocket connections keyed by userID.
type Hub struct {
	mu    sync.RWMutex
	conns map[string]*wsConn // userID → connection
}

func NewHub() *Hub {
	return &Hub{conns: make(map[string]*wsConn)}
}

// IsOnline returns true if the user has an active WebSocket connection.
func (h *Hub) IsOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.conns[userID]
	return ok
}

// Push sends a raw payload to a connected user. Returns false if the user is offline.
func (h *Hub) Push(userID string, payload []byte) bool {
	h.mu.RLock()
	c, ok := h.conns[userID]
	h.mu.RUnlock()
	if !ok {
		return false
	}
	select {
	case c.send <- payload:
		return true
	default:
		return false
	}
}

// ServeWS upgrades the connection and registers the user.
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &wsConn{send: make(chan []byte, 128), conn: conn}

	h.mu.Lock()
	h.conns[userID] = c
	h.mu.Unlock()

	go c.writePump()
	c.readPump()

	h.mu.Lock()
	delete(h.conns, userID)
	h.mu.Unlock()
	close(c.send)
}

func (c *wsConn) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			c.conn.Close()
			return
		}
	}
}

func (c *wsConn) readPump() {
	defer c.conn.Close()
	c.conn.SetReadLimit(512)
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}
