package websocket

import (
	"sync"

	"github.com/ayush/accountability-app/backend/internal/logger"
	"go.uber.org/zap"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients by room ID
	rooms map[string]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations on rooms
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	logger.Info("Creating new WebSocket hub")
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Register adds a new client to the hub
func (h *Hub) Register(client *Client) {
	logger.Info("Registering new client",
		zap.String("room_id", client.RoomID),
		zap.Uint("user_id", client.UserID))
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	logger.Info("Unregistering client",
		zap.String("room_id", client.RoomID),
		zap.Uint("user_id", client.UserID))
	h.unregister <- client
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	logger.Info("Starting WebSocket hub")
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.rooms[client.RoomID]; !ok {
				logger.Info("Creating new room", zap.String("room_id", client.RoomID))
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			logger.Info("Client registered successfully",
				zap.String("room_id", client.RoomID),
				zap.Uint("user_id", client.UserID),
				zap.Int("total_clients_in_room", len(h.rooms[client.RoomID])))
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.rooms[client.RoomID]; ok {
				if _, ok := h.rooms[client.RoomID][client]; ok {
					delete(h.rooms[client.RoomID], client)
					close(client.send)
					logger.Info("Client unregistered",
						zap.String("room_id", client.RoomID),
						zap.Uint("user_id", client.UserID),
						zap.Int("remaining_clients_in_room", len(h.rooms[client.RoomID])))

					// If room is empty, remove it
					if len(h.rooms[client.RoomID]) == 0 {
						delete(h.rooms, client.RoomID)
						logger.Info("Removed empty room", zap.String("room_id", client.RoomID))
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// Broadcast sends a message to all clients in a room
func (h *Hub) Broadcast(roomID string, message []byte) {
	h.mu.RLock()
	if clients, ok := h.rooms[roomID]; ok {
		logger.Debug("Broadcasting message",
			zap.String("room_id", roomID),
			zap.Int("num_clients", len(clients)),
			zap.Int("message_size", len(message)))

		successfulSends := 0
		for client := range clients {
			select {
			case client.send <- message:
				successfulSends++
			default:
				logger.Warn("Failed to send message to client, removing client",
					zap.String("room_id", roomID),
					zap.Uint("user_id", client.UserID))
				close(client.send)
				delete(clients, client)
			}
		}
		logger.Debug("Broadcast complete",
			zap.String("room_id", roomID),
			zap.Int("successful_sends", successfulSends))
	} else {
		logger.Warn("Attempted to broadcast to non-existent room",
			zap.String("room_id", roomID))
	}
	h.mu.RUnlock()
}

// GetClientsInRoom returns the number of clients in a room
func (h *Hub) GetClientsInRoom(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.rooms[roomID]; ok {
		count := len(clients)
		logger.Debug("Retrieved client count for room",
			zap.String("room_id", roomID),
			zap.Int("client_count", count))
		return count
	}
	logger.Debug("Room not found when getting client count",
		zap.String("room_id", roomID))
	return 0
}
