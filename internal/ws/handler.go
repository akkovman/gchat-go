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

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	if !isValidNickname(nickname) {
		sendCloseMessage(conn, InvalidNickname, "invalid nickname")
		return
	}

	// Calls before checking nickname uniqueness
	client := &Client{
		hub:        h,
		conn:       conn,
		nickname:   nickname,
		sendBuffer: make(chan outboundMessage, 256),
	}

	ok, reason := h.addClient(client)
	if !ok {
		switch reason {
		case int(ReasonNicknameIsTaken):
			sendCloseMessage(conn, NicknameIsTaken, "nickname taken")
			return
		case int(ReasonServerIsFull):
			sendCloseMessage(conn, ServerIsFull, "server full")
			return
		}
	}

	h.broadcast(outboundMessage{outboundText, encode("system", "", nickname+" is joined")})

	go client.writePump()
	go client.readPump()
}
