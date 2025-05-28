package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware provides JWT authentication for routes
type AuthMiddleware struct {
	jwtService *JWTService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: NewJWTService(),
	}
}

// extractBearerToken extracts the token from the Authorization header
func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// validateBearerToken validates a JWT token string directly
func (m *AuthMiddleware) validateBearerToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.jwtService.config.SecretKey), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, errors.New("invalid token")
}

// RequireAuth is a middleware that validates JWT tokens
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token *jwt.Token
		var claims jwt.MapClaims
		var err error

		// First try Bearer token authentication
		if bearerToken := extractBearerToken(c); bearerToken != "" {
			token, claims, err = m.validateBearerToken(bearerToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid bearer token"})
				return
			}
		} else {
			// Fall back to cookie-based authentication
			tokenIDStr, err := c.Cookie("jwt-id")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				return
			}

			// Convert token ID to int64
			tokenID, err := strconv.ParseInt(tokenIDStr, 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token ID format"})
				return
			}

			// Validate the token from Redis
			token, claims, err = m.jwtService.ValidateToken(tokenID)
			if err != nil || !token.Valid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				return
			}

			// Store token ID in context for logout
			c.Set("token_id", tokenID)
		}

		// Check token type (only for cookie-based auth)
		if tokenType, ok := claims["type"].(string); ok && tokenType != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
			return
		}

		// Extract user claims
		userClaims, err := m.jwtService.ExtractClaims(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Set user information in context for use in handlers
		c.Set("user_id", userClaims.UserID)
		c.Set("email", userClaims.Email)
		c.Set("role", userClaims.Role)
		c.Set("user_claims", userClaims)

		c.Next()
	}
}

// RequireRole checks if the authenticated user has a specific role
func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This middleware should be used after RequireAuth
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		// Check if user has the required role
		if userRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		c.Next()
	}
}

// GetCurrentUser extracts the user information from the Gin context
func GetCurrentUser(c *gin.Context) (Claims, bool) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return Claims{}, false
	}

	userClaims, ok := claims.(Claims)
	if !ok {
		return Claims{}, false
	}

	return userClaims, true
}

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") //CORS for frontend configuration
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
