package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

/*
Abstract struct containing *websocket.Conn
*/
type Client struct {
	conn     *websocket.Conn
	nickname string
	writeMu  sync.Mutex
}

/*
Send to Client function
*/
func (c *Client) send(data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	return c.conn.WriteMessage(websocket.TextMessage, data)
}
