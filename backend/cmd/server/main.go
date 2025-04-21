package main

import (
	"log"

	"github.com/ayush/accountability-app/backend/docs"
	"github.com/ayush/accountability-app/backend/internal/api"
	"github.com/ayush/accountability-app/backend/internal/config"
	"github.com/ayush/accountability-app/backend/internal/middleware"
	"github.com/ayush/accountability-app/backend/internal/models"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Accountability App API
// @version         1.0
// @description     A video call application with accountability features.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database schemas
	err = db.AutoMigrate(&models.User{}, &models.Call{}, &models.CallParticipant{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize handlers
	userHandler := api.NewUserHandler(db)
	callHandler := api.NewVideoCallHandler(db)
	wsHandler := api.NewWSHandler(cfg.GetWebSocketConfig())

	// Initialize Gin router
	router := gin.Default()

	// Initialize Swagger docs
	docs.SwaggerInfo.BasePath = "/api"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	router.POST("/api/users/register", userHandler.Register)
	router.POST("/api/users/login", userHandler.Login)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/users/:id", userHandler.GetUser)

		// Call routes
		protected.POST("/calls", callHandler.CreateCall)
		protected.POST("/calls/join", callHandler.JoinCall)
		protected.POST("/calls/:room_id/leave", callHandler.LeaveCall)
		protected.GET("/calls", callHandler.ListActiveCalls)

		// WebSocket route
		protected.GET("/ws", wsHandler.HandleWebSocket)
		protected.GET("/rooms/:room_id/participants", wsHandler.GetRoomParticipants)
	}

	// Start the server
	if err := router.Run(cfg.GetServerAddress()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
