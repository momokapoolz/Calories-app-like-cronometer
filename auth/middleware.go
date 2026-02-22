package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
)

// AuthMiddleware provides JWT authentication for routes
type AuthMiddleware struct {
	jwtService *JWTService
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{jwtService: NewJWTService()}
}

// extractBearerToken reads the JWT string from the Authorization: Bearer header
func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

// RequireAuth validates a JWT access token from the Bearer header or the
// access_token HttpOnly cookie. On success it sets user_id, email, role, and
// user_claims in the Gin context for downstream handlers.
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Bearer header takes priority (API clients, Postman, mobile apps)
		if bearer := extractBearerToken(c); bearer != "" {
			tokenString = bearer
		} else {
			cookie, err := c.Cookie(AccessTokenCookie)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status":  "error",
					"message": "Authentication required",
				})
				return
			}
			tokenString = cookie
		}

		_, claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid or expired token",
			})
			return
		}

		// Reject non-access tokens (e.g. someone accidentally sending a refresh token)
		if tokenType, ok := claims["type"].(string); !ok || tokenType != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid token type",
			})
			return
		}

		userClaims, err := m.jwtService.ExtractClaims(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid token claims",
			})
			return
		}

		c.Set("user_id", userClaims.UserID)
		c.Set("email", userClaims.Email)
		c.Set("role", userClaims.Role)
		c.Set("user_claims", userClaims)

		c.Next()
	}
}

// RequireRole checks that the authenticated user has the given role.
// Must be chained after RequireAuth.
func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authentication required",
			})
			return
		}
		if userRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Insufficient permissions",
			})
			return
		}
		c.Next()
	}
}

// GetCurrentUser extracts the typed Claims from the Gin context (set by RequireAuth)
func GetCurrentUser(c *gin.Context) (Claims, bool) {
	val, exists := c.Get("user_claims")
	if !exists {
		return Claims{}, false
	}
	userClaims, ok := val.(Claims)
	return userClaims, ok
}

// CORSMiddleware sets CORS headers. The allowed origin is read from the
// CORS_ORIGIN env variable; falls back to http://localhost:3000 for development.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := os.Getenv("CORS_ORIGIN")
		if origin == "" {
			origin = "http://localhost:3000"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
