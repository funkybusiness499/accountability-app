package api

import (
	"net/http"

	"github.com/ayush/accountability-app/backend/internal/auth"
	"github.com/ayush/accountability-app/backend/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

// UserResponse represents the user information in responses
type UserResponse struct {
	ID       uint   `json:"id" example:"1"`
	Email    string `json:"email" example:"john@example.com"`
	Username string `json:"username" example:"johndoe"`
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User registration details"
// @Success 201 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to hash password"})
		return
	}

	user := models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, LoginResponse{
		Token: token,
		User: UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	})
}

// GetUser godoc
// @Summary Get user details
// @Description Get user information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 404 {object} ErrorResponse
// @Security Bearer
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	})
}
