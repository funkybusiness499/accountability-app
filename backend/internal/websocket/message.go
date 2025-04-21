package websocket

import (
	"encoding/json"
	"time"

	"github.com/ayush/accountability-app/backend/internal/logger"
	"go.uber.org/zap"
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
	msg := &Message{
		Type:      msgType,
		Data:      data,
		RoomID:    roomID,
		UserID:    userID,
		Timestamp: time.Now(),
	}

	logger.Debug("Created new message",
		zap.String("type", string(msgType)),
		zap.String("room_id", roomID),
		zap.Uint("user_id", userID))

	return msg
}

// Marshal converts the message to JSON bytes
func (m *Message) Marshal() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		logger.Error("Failed to marshal message",
			zap.Error(err),
			zap.String("type", string(m.Type)),
			zap.String("room_id", m.RoomID),
			zap.Uint("user_id", m.UserID))
		return nil, err
	}

	logger.Debug("Marshaled message",
		zap.String("type", string(m.Type)),
		zap.String("room_id", m.RoomID),
		zap.Uint("user_id", m.UserID),
		zap.Int("bytes_length", len(bytes)))

	return bytes, nil
}

// UnmarshalMessage converts JSON bytes to a Message
func UnmarshalMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		logger.Error("Failed to unmarshal message",
			zap.Error(err),
			zap.Binary("raw_data", data))
		return nil, err
	}

	logger.Debug("Unmarshaled message",
		zap.String("type", string(msg.Type)),
		zap.String("room_id", msg.RoomID),
		zap.Uint("user_id", msg.UserID))

	return &msg, nil
}
