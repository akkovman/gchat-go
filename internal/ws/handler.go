package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) ServeWS(ctx *gin.Context) {
	nickname := ctx.Query("nickname")
	if nickname == "" {
		ctx.String(http.StatusBadRequest, "Nickname is required")
		return
	}

	// TODO: Add a unique nickname system later

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:        h,
		conn:       conn,
		nickname:   nickname,
		sendBuffer: make(chan []byte, 256),
	}
	h.addClient(client)
	h.broadcast(encode("system", "", nickname+" is joined"))

	go client.writePump()
	go client.readPump()
}
