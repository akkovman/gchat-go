package ws

import (
	"encoding/json"
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
		hub:      h,
		conn:     conn,
		nickname: nickname,
	}
	h.addClient(client)
	h.broadcast(encode("system", "", nickname+" is joined"))

	defer func() {
		h.deleteClient(client)
		h.broadcast(encode("system", "", nickname+" is left"))

		conn.Close()
	}()

	for {
		_, raw, err := conn.ReadMessage()
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

		h.broadcast(encode("message", nickname, in.Text))
	}
}
