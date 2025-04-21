package websocket

import (
	"sync"
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
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Register adds a new client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.rooms[client.RoomID]; !ok {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.rooms[client.RoomID]; ok {
				if _, ok := h.rooms[client.RoomID][client]; ok {
					delete(h.rooms[client.RoomID], client)
					close(client.send)
					// If room is empty, remove it
					if len(h.rooms[client.RoomID]) == 0 {
						delete(h.rooms, client.RoomID)
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
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
	}
	h.mu.RUnlock()
}

// GetClientsInRoom returns the number of clients in a room
func (h *Hub) GetClientsInRoom(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.rooms[roomID]; ok {
		return len(clients)
	}
	return 0
}
