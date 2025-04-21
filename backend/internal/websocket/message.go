package websocket

import (
	"encoding/json"
	"time"
)

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Message types
	MessageTypeChat     MessageType = "chat"
	MessageTypePresence MessageType = "presence"
	MessageTypeSystem   MessageType = "system"
)

// Message represents a structured WebSocket message
type Message struct {
	Type      MessageType `json:"type"`
	Data      interface{} `json:"data"`
	RoomID    string      `json:"room_id"`
	UserID    uint        `json:"user_id"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewMessage creates a new message with the current timestamp
func NewMessage(msgType MessageType, data interface{}, roomID string, userID uint) *Message {
	return &Message{
		Type:      msgType,
		Data:      data,
		RoomID:    roomID,
		UserID:    userID,
		Timestamp: time.Now(),
	}
}

// Marshal converts the message to JSON bytes
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalMessage converts JSON bytes to a Message
func UnmarshalMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
