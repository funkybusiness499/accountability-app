package models

import (
	"time"
)

// Call represents a video call session
type Call struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   uint      `json:"creator_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CallCreate represents the request to create a new call
type CallCreate struct {
	Title       string `json:"title" binding:"required" example:"Team Meeting"`
	Description string `json:"description" example:"Weekly team sync meeting"`
	CreatorID   uint   `json:"creator_id" binding:"required"`
}

// CallJoin represents the request to join a call
type CallJoin struct {
	CallID uint `json:"call_id" binding:"required" example:"1"`
	UserID uint `json:"user_id" binding:"required" example:"2"`
}

// CallParticipant represents a user participating in a call
type CallParticipant struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CallID    uint      `json:"call_id"`
	UserID    uint      `json:"user_id"`
	JoinedAt  time.Time `json:"joined_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for the Call model
func (Call) TableName() string {
	return "calls"
}

// TableName specifies the table name for the CallParticipant model
func (CallParticipant) TableName() string {
	return "call_participants"
}
