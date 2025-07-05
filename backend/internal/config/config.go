package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Auth     AuthConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL        string
	MaxRetries int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	AdminEmail    string
	AdminPassword string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config := &Config{
		Database: DatabaseConfig{
			URL:        getEnv("DATABASE_URL", ""),
			MaxRetries: getEnvAsInt("DB_MAX_RETRIES", 10),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Auth: AuthConfig{
			AdminEmail:    getEnv("ADMIN_EMAIL", ""),
			AdminPassword: getEnv("ADMIN_PASSWORD", ""),
		},
	}

	// Validate required fields
	if config.Database.URL == "" {
		return nil, &ConfigError{Field: "DATABASE_URL", Message: "is required"}
	}

	return config, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// ConfigError represents a configuration error
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Field + " " + e.Message
}