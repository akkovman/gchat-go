package ws

import (
	"encoding/json"
	"time"
)

/*
Abstract struct with json params for frontend
*/
type outMessage struct {
	Type     string `json:"type"`
	Nickname string `json:"nickname,omitempty"`
	Text     string `json:"text"`
	Ts       int64  `json:"ts"`
}

func encode(msgType, nickname, text string) []byte {
	data, _ := json.Marshal(outMessage{
		Type: msgType, Nickname: nickname, Text: text, Ts: time.Now().UnixMilli(),
	})

	return data
}

/*
Enum types of message
*/
type outboundKind int

const (
	outboundText outboundKind = iota
	outboundClose
)

/*
Struct which contains type of message and data
*/
type outboundMessage struct {
	kind outboundKind
	data []byte
}
