package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewDatabaseConfig creates a new database configuration
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:   getEnvOrDefault("DB_NAME", "videocall"),
	}
}

// Connect establishes a connection to the database
func (c *DatabaseConfig) Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.DBName, c.Port)
	
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// Helper function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 