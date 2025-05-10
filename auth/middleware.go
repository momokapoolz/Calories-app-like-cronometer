package auth

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// RequireAuth is a middleware that validates JWT tokens
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token ID from cookie
		tokenIDStr, err := c.Cookie("jwt-id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "JWT token ID cookie required"})
			return
		}

		// Convert token ID to int64
		tokenID, err := strconv.ParseInt(tokenIDStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token ID format"})
			return
		}

		// Validate the token
		token, claims, err := m.jwtService.ValidateToken(tokenID)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Check token type
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
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
		c.Set("token_id", tokenID)

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