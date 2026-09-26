package ws

import (
	"log"
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
	delete(h.clients, c)
}

/*
Send to everyone function
*/
func (h *Hub) broadcast(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for conn := range h.clients {
		if err := conn.send(message); err != nil {
			log.Printf("Broadcast error: %v", err)
		}
	}
}
