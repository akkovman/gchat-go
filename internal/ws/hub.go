package ws

import (
	"gchat/internal/env"
	"gchat/internal/ip"
	"sync"

	"github.com/gorilla/websocket"
)

type RejectReason int

const (
	ReasonNothing RejectReason = iota
	ReasonNicknameIsTaken
	ReasonServerIsFull
)

type Hub struct {
	mu sync.RWMutex

	clients map[*Client]bool
	WG      sync.WaitGroup

	BanList *ip.BanList
}

func NewHub(bl *ip.BanList) *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		BanList: bl,
	}
}

func (h *Hub) addClient(c *Client) (bool, int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// BEGIN POLICY
	for exists := range h.clients {
		if exists.nickname == c.nickname {
			return false, int(ReasonNicknameIsTaken)
		}
	}
	// END POLICY

	// BEGIN SERVER FULL
	if len(h.clients) >= env.GetEnvAsInt("MAX_CLIENT_CONNECTIONS", 1024) {
		return false, int(ReasonServerIsFull)
	}
	// END SERVER FULL

	h.clients[c] = true
	return true, int(ReasonNothing)
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
func (h *Hub) CloseAll() {
	payload := websocket.FormatCloseMessage(1001, "Server is shutting down")
	h.broadcast(outboundMessage{outboundClose, payload})
}
