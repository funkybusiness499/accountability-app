package websocket

import (
	"time"

	"github.com/ayush/accountability-app/backend/internal/logger"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 32768
)

// Client represents a websocket connection in the hub
type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// The hub instance
	hub *Hub

	// Buffered channel of outbound messages
	send chan []byte

	// Room ID this client belongs to
	RoomID string

	// User ID associated with this client
	UserID uint
}

// NewClient creates a new client instance
func NewClient(hub *Hub, conn *websocket.Conn, roomID string, userID uint) *Client {
	logger.Info("Creating new WebSocket client",
		zap.String("room_id", roomID),
		zap.Uint("user_id", userID))
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		RoomID: roomID,
		UserID: userID,
	}
}

// ReadPump pumps messages from the websocket connection to the hub
func (c *Client) ReadPump() {
	logger.Info("Starting client read pump",
		zap.String("room_id", c.RoomID),
		zap.Uint("user_id", c.UserID))

	defer func() {
		logger.Info("Closing client read pump",
			zap.String("room_id", c.RoomID),
			zap.Uint("user_id", c.UserID))
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		logger.Debug("Received pong from client",
			zap.String("room_id", c.RoomID),
			zap.Uint("user_id", c.UserID))
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("Unexpected WebSocket close",
					zap.Error(err),
					zap.String("room_id", c.RoomID),
					zap.Uint("user_id", c.UserID))
			} else {
				logger.Info("WebSocket connection closed normally",
					zap.String("room_id", c.RoomID),
					zap.Uint("user_id", c.UserID))
			}
			break
		}

		// Parse the incoming message
		msg, err := UnmarshalMessage(rawMessage)
		if err != nil {
			logger.Error("Failed to parse incoming message",
				zap.Error(err),
				zap.String("room_id", c.RoomID),
				zap.Uint("user_id", c.UserID),
				zap.Binary("raw_message", rawMessage))
			continue
		}

		// Ensure the message is for this room
		if msg.RoomID != c.RoomID {
			logger.Warn("Received message for wrong room",
				zap.String("expected_room", c.RoomID),
				zap.String("received_room", msg.RoomID),
				zap.Uint("user_id", c.UserID))
			continue
		}

		// Set the correct UserID from the client
		msg.UserID = c.UserID

		// Marshal the message back to JSON
		messageBytes, err := msg.Marshal()
		if err != nil {
			logger.Error("Failed to marshal message for broadcast",
				zap.Error(err),
				zap.String("room_id", c.RoomID),
				zap.Uint("user_id", c.UserID))
			continue
		}

		logger.Debug("Broadcasting message from client",
			zap.String("room_id", c.RoomID),
			zap.Uint("user_id", c.UserID),
			zap.Int("message_size", len(messageBytes)))

		// Broadcast the message to all clients in the same room
		c.hub.Broadcast(c.RoomID, messageBytes)
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
	logger.Info("Starting client write pump",
		zap.String("room_id", c.RoomID),
		zap.Uint("user_id", c.UserID))

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		logger.Info("Closing client write pump",
			zap.String("room_id", c.RoomID),
			zap.Uint("user_id", c.UserID))
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				logger.Info("Hub closed client send channel",
					zap.String("room_id", c.RoomID),
					zap.Uint("user_id", c.UserID))
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Error("Failed to get next writer",
					zap.Error(err),
					zap.String("room_id", c.RoomID),
					zap.Uint("user_id", c.UserID))
				return
			}

			logger.Debug("Writing message to client",
				zap.String("room_id", c.RoomID),
				zap.Uint("user_id", c.UserID),
				zap.Int("message_size", len(message)))

			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			if n > 0 {
				logger.Debug("Processing queued messages",
					zap.String("room_id", c.RoomID),
					zap.Uint("user_id", c.UserID),
					zap.Int("queued_messages", n))
			}

			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				logger.Error("Error closing writer",
					zap.Error(err),
					zap.String("room_id", c.RoomID),
					zap.Uint("user_id", c.UserID))
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			logger.Debug("Sending ping to client",
				zap.String("room_id", c.RoomID),
				zap.Uint("user_id", c.UserID))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Error("Failed to send ping",
					zap.Error(err),
					zap.String("room_id", c.RoomID),
					zap.Uint("user_id", c.UserID))
				return
			}
		}
	}
}
