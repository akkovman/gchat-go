package ws

import "github.com/gorilla/websocket"

/*
Send to peer CloseMessage with close code 4001 (nickname is taken)
*/
func rejectNicknameTaken(conn *websocket.Conn) {
	payload := websocket.FormatCloseMessage(4001, "nickname taken")
	conn.WriteMessage(websocket.CloseMessage, payload)
	conn.Close()
}
