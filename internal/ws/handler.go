package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// BEGIN CHECK ORIGIN

// Check origin should through nginx proxy
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// END CHECK OIRIGN

func (h *Hub) ServeWS(ctx *gin.Context) {
	nickname := ctx.Query("nickname")

	ip := ctx.ClientIP()
	if h.BanList.Check(ip) {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

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

	h.WG.Add(2)
	go client.writePump()
	go client.readPump()
}
