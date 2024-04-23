package models

import (
	"bytes"
	"fmt"
)

type ActionType string

const (
	DIRECT    ActionType = "direct"
	BROADCAST ActionType = "broadcast"
	PING      ActionType = "ping"
	PONG      ActionType = "pong"
)

type WsMessage struct {
	ActionType ActionType `json:"action_type"`
	Body       string     `json:"body"`
	Data       string     `json:"data,omitempty"`
	UserId     string     `json:"user_id,omitempty"`
}

func NewMessage(actionType ActionType, userId, body, data string) WsMessage {
	return WsMessage{
		ActionType: actionType,
		UserId:     userId,
		Body:       body,
		Data:       data,
	}
}

func (w *WsMessage) Show() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("{ActionType: %s, ", w.ActionType))
	buf.WriteString(fmt.Sprintf("UserId: %s, ", w.UserId))
	buf.WriteString(fmt.Sprintf("Body: %s, ", w.Body))
	buf.WriteString(fmt.Sprintf("Data: %s}\n", w.Data))

	return buf.String()
}
