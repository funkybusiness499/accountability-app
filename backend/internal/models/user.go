package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Password  string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRegistration represents the request to register a new user
type UserRegistration struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"secretpassword"`
}

// UserLogin represents the request to login
type UserLogin struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"secretpassword"`
}

// CheckPassword verifies if the provided password matches the hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// GenerateToken generates a JWT token for the user
func (u *User) GenerateToken() (string, error) {
	// TODO: Implement JWT token generation
	return "", nil
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
