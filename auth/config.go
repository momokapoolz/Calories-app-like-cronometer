package auth

import (
	"os"
	"time"
)

// Config holds all JWT and cookie configuration
type Config struct {
	SecretKey     string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	Issuer        string
	CookieSecure  bool   // set to true in production (HTTPS)
	CookieDomain  string // empty = current host
}

// GetConfig returns configuration loaded from environment variables with safe defaults
func GetConfig() Config {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key-change-in-production"
	}

	// Set COOKIE_SECURE=true in production; default false for local HTTP dev
	secure := os.Getenv("COOKIE_SECURE") == "true"
	domain := os.Getenv("COOKIE_DOMAIN")

	return Config{
		SecretKey:     secretKey,
		TokenExpiry:   time.Minute * 15,   // Access token: 15 minutes (industry standard)
		RefreshExpiry: time.Hour * 24 * 7, // Refresh token: 7 days
		Issuer:        "caloriesapp",
		CookieSecure:  secure,
		CookieDomain:  domain,
	}
}
