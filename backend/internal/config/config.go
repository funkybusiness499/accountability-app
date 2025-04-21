package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`

	JWT struct {
		Secret          string `yaml:"secret"`
		ExpirationHours int    `yaml:"expiration_hours"`
	} `yaml:"jwt"`

	WebSocket struct {
		AllowedOrigins []string `yaml:"allowed_origins"`
	} `yaml:"websocket"`

	Server struct {
		Port        int    `yaml:"port"`
		Environment string `yaml:"environment"`
	} `yaml:"server"`
}

var globalConfig *Config

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	// Try to find config file in different locations
	configPaths := []string{
		"config.yaml",
		"../config.yaml",
		"../../config.yaml",
	}

	var configData []byte
	var err error

	for _, path := range configPaths {
		configData, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(configData, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Override with environment variables if they exist
	if envHost := os.Getenv("DB_HOST"); envHost != "" {
		config.Database.Host = envHost
	}
	if envPort := os.Getenv("DB_PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", &config.Database.Port)
	}
	if envUser := os.Getenv("DB_USER"); envUser != "" {
		config.Database.User = envUser
	}
	if envPassword := os.Getenv("DB_PASSWORD"); envPassword != "" {
		config.Database.Password = envPassword
	}
	if envName := os.Getenv("DB_NAME"); envName != "" {
		config.Database.Name = envName
	}
	if envSecret := os.Getenv("JWT_SECRET"); envSecret != "" {
		config.JWT.Secret = envSecret
	}

	globalConfig = config
	return config, nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Name,
		c.Database.SSLMode,
	)

	if c.Database.Password != "" {
		dsn += fmt.Sprintf(" password=%s", c.Database.Password)
	}

	return dsn
}

// GetWebSocketConfig returns the WebSocket configuration
func (c *Config) GetWebSocketConfig() *WebSocketConfig {
	return &WebSocketConfig{
		AllowedOrigins: c.WebSocket.AllowedOrigins,
	}
}

// GetServerAddress returns the server address with port
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}
