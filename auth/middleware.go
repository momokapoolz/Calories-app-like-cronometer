package auth

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
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

// isUUID checks if a string is a valid UUID format
func isUUID(s string) bool {
	// UUID regex pattern (supports both v4 and other UUID formats)
	uuidPattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	matched, _ := regexp.MatchString(uuidPattern, s)
	return matched
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

// validateBearerTokenOrID validates either a JWT token string or a token ID (UUID)
func (m *AuthMiddleware) validateBearerTokenOrID(tokenString string) (*jwt.Token, jwt.MapClaims, string, error) {
	// Check if the token is a UUID (token ID)
	if isUUID(tokenString) {
		// It's a token ID, retrieve the actual JWT from Redis
		token, claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			return nil, nil, "", fmt.Errorf("failed to validate token ID: %v", err)
		}
		return token, claims, tokenString, nil
	} else {
		// It's a direct JWT token, validate it directly
		token, claims, err := m.validateBearerToken(tokenString)
		if err != nil {
			return nil, nil, "", err
		}
		return token, claims, "", nil
	}
}

// RequireAuth is a middleware that validates JWT tokens
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token *jwt.Token
		var claims jwt.MapClaims
		var tokenID string
		var err error

		// First try Bearer token authentication
		if bearerToken := extractBearerToken(c); bearerToken != "" {
			token, claims, tokenID, err = m.validateBearerTokenOrID(bearerToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid bearer token"})
				return
			}

			// Store token ID in context for logout if it was a UUID
			if tokenID != "" {
				c.Set("token_id", tokenID)
			}
		} else {
			// Fall back to cookie-based authentication
			tokenIDStr, err := c.Cookie("jwt-id")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				return
			}

			// Use the token ID directly (now a string)
			tokenID = tokenIDStr

			// Validate the token from Redis
			token, claims, err = m.jwtService.ValidateToken(tokenID)
			if err != nil || !token.Valid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				return
			}

			// Store token ID in context for logout
			c.Set("token_id", tokenID)
		}

		// Check token type (only for cookie-based auth or token ID based auth)
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
