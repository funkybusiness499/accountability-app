package config

// WebSocketConfig holds configuration for WebSocket connections
type WebSocketConfig struct {
	// AllowedOrigins is a list of origins allowed to connect to the WebSocket server
	AllowedOrigins []string
}

// DefaultWebSocketConfig returns the default WebSocket configuration
func DefaultWebSocketConfig() *WebSocketConfig {
	return &WebSocketConfig{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
	}
}
