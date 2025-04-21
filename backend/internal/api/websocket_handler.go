// Package api provides HTTP and WebSocket API handlers for the application.
package api

import (
	"net/http"
	"strconv"

	"github.com/ayush/accountability-app/backend/internal/config"
	"github.com/ayush/accountability-app/backend/internal/logger"
	ws "github.com/ayush/accountability-app/backend/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// WSHandler handles WebSocket connections
type WSHandler struct {
	hub    *ws.Hub
	config *config.WebSocketConfig
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Note: In production, implement proper origin checking
		logger.Debug("Checking WebSocket origin",
			zap.String("origin", r.Header.Get("Origin")))
		return true
	},
}

// NewWSHandler creates a new WebSocket handler
func NewWSHandler(config *config.WebSocketConfig) *WSHandler {
	logger.Info("Creating new WebSocket handler",
		zap.Strings("allowed_origins", config.AllowedOrigins))

	hub := ws.NewHub()
	go hub.Run()
	return &WSHandler{
		hub:    hub,
		config: config,
	}
}

// HandleWebSocket godoc
// @Summary Connect to WebSocket
// @Description Establish a WebSocket connection for real-time communication
// @Tags websocket
// @Accept json
// @Produce json
// @Param room_id query string true "Room ID"
// @Param user_id query string true "User ID"
// @Success 101 {string} string "Switching Protocols to websocket"
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security Bearer
// @Router /ws [get]
func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	roomID := c.Query("room_id")
	if roomID == "" {
		logger.Warn("WebSocket connection attempt without room_id")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "room_id is required"})
		return
	}

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		logger.Warn("WebSocket connection attempt without user_id",
			zap.String("room_id", roomID))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "user_id is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logger.Warn("Invalid user_id format in WebSocket connection attempt",
			zap.String("room_id", roomID),
			zap.String("user_id_str", userIDStr),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user_id format"})
		return
	}

	logger.Info("WebSocket connection attempt",
		zap.String("room_id", roomID),
		zap.Uint64("user_id", userID),
		zap.String("remote_addr", c.Request.RemoteAddr),
		zap.String("origin", c.Request.Header.Get("Origin")))

	// Check origin if configured
	if len(h.config.AllowedOrigins) > 0 {
		origin := c.Request.Header.Get("Origin")
		allowed := false
		for _, allowedOrigin := range h.config.AllowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}
		if !allowed {
			logger.Warn("WebSocket connection rejected due to unauthorized origin",
				zap.String("origin", origin),
				zap.String("room_id", roomID),
				zap.Uint64("user_id", userID),
				zap.Strings("allowed_origins", h.config.AllowedOrigins))
			c.JSON(http.StatusForbidden, ErrorResponse{Error: "origin not allowed"})
			return
		}
	}

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Failed to upgrade WebSocket connection",
			zap.Error(err),
			zap.String("room_id", roomID),
			zap.Uint64("user_id", userID))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to upgrade connection"})
		return
	}

	logger.Info("WebSocket connection established",
		zap.String("room_id", roomID),
		zap.Uint64("user_id", userID),
		zap.String("remote_addr", c.Request.RemoteAddr))

	client := ws.NewClient(h.hub, conn, roomID, uint(userID))
	h.hub.Register(client)

	// Start client message pumps
	go client.WritePump()
	go client.ReadPump()
}

// GetRoomParticipants godoc
// @Summary Get room participants
// @Description Get the number of participants in a specific room
// @Tags websocket
// @Accept json
// @Produce json
// @Param room_id path string true "Room ID"
// @Success 200 {object} WSParticipantsResponse
// @Failure 400 {object} ErrorResponse
// @Security Bearer
// @Router /rooms/{room_id}/participants [get]
func (h *WSHandler) GetRoomParticipants(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		logger.Warn("Room participants request without room_id")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "room_id is required"})
		return
	}

	count := h.hub.GetClientsInRoom(roomID)
	logger.Debug("Retrieved room participants count",
		zap.String("room_id", roomID),
		zap.Int("count", count))

	c.JSON(http.StatusOK, WSParticipantsResponse{Count: count})
}

// WSParticipantsResponse represents the response for room participants count
type WSParticipantsResponse struct {
	Count int `json:"count"`
}
