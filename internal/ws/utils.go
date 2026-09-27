package ws

import "github.com/gorilla/websocket"

const (
	NicknameIsTaken = 4001
	ServerIsFull    = 4002
)

/*
Send to peer CloseMessage with close code
*/
func sendCloseMessage(conn *websocket.Conn, code int, message string) {
	payload := websocket.FormatCloseMessage(code, message)
	conn.WriteMessage(websocket.CloseMessage, payload)
	conn.Close()
}
