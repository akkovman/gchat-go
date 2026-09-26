package ws

import (
	"sync"
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

func (h *Hub) addClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[c] = true
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
func (h *Hub) broadcast(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for conn := range h.clients {
		select {
		case conn.sendBuffer <- message:
		default:
			close(conn.sendBuffer)
			delete(h.clients, conn)
		}
	}
}
