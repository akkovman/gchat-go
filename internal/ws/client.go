package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write message to peer
	writeWait = 10 * time.Second

	// Time allowed to read a pong message from peer
	pongWait = 50 * time.Second

	// Send ping to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

/*
Abstract struct containing *websocket.Conn
*/
type Client struct {
	hub        *Hub
	conn       *websocket.Conn
	nickname   string
	writeMu    sync.Mutex
	sendBuffer chan []byte // Buffered channel for i/o
}

/*
Send to Client function
*/
func (c *Client) send(data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	return c.conn.WriteMessage(websocket.TextMessage, data)
}

/*
Pumps message from the websocket connection to the hub
Runs in a goroutine
*/
func (c *Client) readPump() {
	defer func() {
		c.hub.deleteClient(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var in struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}

		if json.Unmarshal(raw, &in) != nil || in.Type != "message" || in.Text == "" {
			continue
		}

		c.hub.broadcast(encode("message", c.nickname, in.Text))
	}
}

/*
Pumps message from the hub to the websocket connection
*/
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()

		c.hub.broadcast(encode("system", "", c.nickname+" is left"))
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.sendBuffer:
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
