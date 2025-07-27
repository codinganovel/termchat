package protocol

import "time"

type MessageType string

const (
	MessageTypeHello   MessageType = "hello"
	MessageTypeWelcome MessageType = "welcome"
	MessageTypeReady   MessageType = "ready"
	MessageTypeText    MessageType = "text"
	MessageTypePing    MessageType = "ping"
	MessageTypePong    MessageType = "pong"
	MessageTypeLeave   MessageType = "leave"
	MessageTypeError   MessageType = "error"
)

type Message struct {
	Type      MessageType `json:"type"`
	Content   string      `json:"content,omitempty"`
	SessionID string      `json:"session_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

func NewMessage(msgType MessageType, content string) *Message {
	return &Message{
		Type:      msgType,
		Content:   content,
		Timestamp: time.Now().UnixMilli(),
	}
}

func NewHandshakeMessage(msgType MessageType, sessionID string) *Message {
	return &Message{
		Type:      msgType,
		SessionID: sessionID,
		Timestamp: time.Now().UnixMilli(),
	}
}