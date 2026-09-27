package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
	}
}

func (h *Hub) addClient(c *Client) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	// BEGIN POLICY
	for exists := range h.clients {
		if exists.nickname == c.nickname {
			return false
		}
	}
	// END POLICY

	h.clients[c] = true
	return true
}

func (h *Hub) deleteClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[c]; exists {
		delete(h.clients, c)
		close(c.sendBuffer)
	}
}

/*
Send to everyone function
*/
func (h *Hub) broadcast(outMsg outboundMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for conn := range h.clients {
		select {
		case conn.sendBuffer <- outMsg:
		default:
			close(conn.sendBuffer)
			delete(h.clients, conn)
		}
	}
}

/*
Close all websocket connections function
*/
func (h *Hub) closeAll() {
	payload := websocket.FormatCloseMessage(1001, "Server is shutting down")
	h.broadcast(outboundMessage{outboundClose, payload})
}
