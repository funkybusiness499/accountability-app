package api

import (
	"net/http"
	"time"

	"github.com/ayush/accountability-app/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// VideoCallHandler handles video call-related HTTP endpoints
type VideoCallHandler struct {
	db *gorm.DB
}

// NewVideoCallHandler creates a new video call handler
func NewVideoCallHandler(db *gorm.DB) *VideoCallHandler {
	return &VideoCallHandler{db: db}
}

// CreateCall godoc
// @Summary Create a new call
// @Description Create a new video call session
// @Tags calls
// @Accept json
// @Produce json
// @Param call body models.CallCreate true "Call details"
// @Success 201 {object} models.Call
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security Bearer
// @Router /calls [post]
func (h *VideoCallHandler) CreateCall(c *gin.Context) {
	var input models.CallCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	call := models.Call{
		Title:       input.Title,
		Description: input.Description,
		CreatorID:   input.CreatorID,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.db.Create(&call).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create call"})
		return
	}

	c.JSON(http.StatusCreated, call)
}

// JoinCall godoc
// @Summary Join an existing call
// @Description Join an active video call session
// @Tags calls
// @Accept json
// @Produce json
// @Param join body models.CallJoin true "Join call details"
// @Success 200 {object} models.CallParticipant
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security Bearer
// @Router /calls/join [post]
func (h *VideoCallHandler) JoinCall(c *gin.Context) {
	var input models.CallJoin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var call models.Call
	if err := h.db.First(&call, input.CallID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Call not found"})
		return
	}

	if call.Status != "active" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Call is not active"})
		return
	}

	participant := models.CallParticipant{
		CallID:    input.CallID,
		UserID:    input.UserID,
		JoinedAt:  time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.db.Create(&participant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to join call"})
		return
	}

	c.JSON(http.StatusOK, participant)
}

// LeaveCall godoc
// @Summary Leave a call
// @Description Leave a video call session
// @Tags calls
// @Accept json
// @Produce json
// @Param room_id path string true "Room ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security Bearer
// @Router /calls/{room_id}/leave [post]
func (h *VideoCallHandler) LeaveCall(c *gin.Context) {
	roomID := c.Param("room_id")
	userID := c.GetString("user_id")

	if err := h.db.Where("call_id = ? AND user_id = ?", roomID, userID).Delete(&models.CallParticipant{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to leave call"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Successfully left the call"})
}

// ListActiveCalls godoc
// @Summary List all active calls
// @Description Get a list of all active video call sessions
// @Tags calls
// @Accept json
// @Produce json
// @Success 200 {array} models.Call
// @Failure 500 {object} ErrorResponse
// @Security Bearer
// @Router /calls [get]
func (h *VideoCallHandler) ListActiveCalls(c *gin.Context) {
	var calls []models.Call
	if err := h.db.Where("status = ?", "active").Find(&calls).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch calls"})
		return
	}

	c.JSON(http.StatusOK, calls)
}
