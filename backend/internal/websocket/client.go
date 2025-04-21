package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
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
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Parse the incoming message
		msg, err := UnmarshalMessage(rawMessage)
		if err != nil {
			log.Printf("error parsing message: %v", err)
			continue
		}

		// Ensure the message is for this room
		if msg.RoomID != c.RoomID {
			log.Printf("warning: received message for wrong room: got %s, expected %s", msg.RoomID, c.RoomID)
			continue
		}

		// Set the correct UserID from the client
		msg.UserID = c.UserID

		// Marshal the message back to JSON
		messageBytes, err := msg.Marshal()
		if err != nil {
			log.Printf("error marshaling message: %v", err)
			continue
		}

		// Broadcast the message to all clients in the same room
		c.hub.Broadcast(c.RoomID, messageBytes)
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
