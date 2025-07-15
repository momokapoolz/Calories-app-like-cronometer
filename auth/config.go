package auth

import (
	"os"
	"time"
)

// Config holds all JWT configuration
type Config struct {
	SecretKey     string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

// GetConfig returns JWT configuration with sensible defaults
func GetConfig() Config {
	// In production, these values should be loaded from environment variables
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key-change-in-production" // Default for development only
	}

	return Config{
		SecretKey:     secretKey,
		TokenExpiry:   time.Hour * 10,     // Access tokens valid for 10 hour
		RefreshExpiry: time.Hour * 24 * 7, // Refresh tokens valid for 1 week
		Issuer:        "caloriesapp",
	}
}
