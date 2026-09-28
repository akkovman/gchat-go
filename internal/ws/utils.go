package ws

import (
	"regexp"

	"github.com/gorilla/websocket"
)

const (
	NicknameIsTaken = 4001
	ServerIsFull    = 4002
	InvalidNickname = 4003
)

/*
Send to peer CloseMessage with close code
*/
func sendCloseMessage(conn *websocket.Conn, code int, message string) {
	payload := websocket.FormatCloseMessage(code, message)
	conn.WriteMessage(websocket.CloseMessage, payload)
	conn.Close()
}

// BEGIN VALIDATOR

/*
Only latin characters, numbers and underscores are allowed
3 to 20 characters
*/
var nicknameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

/*
Function validate nickname
*/
func isValidNickname(nickname string) bool {
	return nicknameRegex.MatchString(nickname)
}

// END VALIDATOR
