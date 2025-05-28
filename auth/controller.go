package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/momokapoolz/caloriesapp/user/models"
	"github.com/momokapoolz/caloriesapp/user/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthController handles authentication endpoints
type AuthController struct {
	jwtService *JWTService
	userRepo   *repository.UserRepository
}

// LoginResponse represents the response for both token-based and cookie-based auth
type LoginResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// TokenResponse represents the response for token-based authentication
type TokenResponse struct {
	User      interface{} `json:"user"`
	Token     string      `json:"token"`
	ExpiresIn int64       `json:"expires_in"`
}

// CookieResponse represents the response for cookie-based authentication
type CookieResponse struct {
	User      interface{} `json:"user"`
	ExpiresIn int64       `json:"expires_in"`
}

// NewAuthController creates a new auth controller
func NewAuthController(userRepo *repository.UserRepository) *AuthController {
	return &AuthController{
		jwtService: NewJWTService(),
		userRepo:   userRepo,
	}
}

// Login authenticates a user and generates JWT tokens
func (c *AuthController) Login(ctx *gin.Context) {
	var request LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}

	// Find user by email
	user, err := c.userRepo.FindByEmail(request.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token pair
	tokenPair, err := c.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Check if client wants token-based or cookie-based auth
	authType := ctx.GetHeader("X-Auth-Type")
	if authType == "token" {
		// Get the actual token for token-based auth
		accessToken, err := GetToken(tokenPair.AccessTokenID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve token"})
			return
		}

		response := LoginResponse{
			Status:  "success",
			Message: "Login successful",
			Data: TokenResponse{
				User: gin.H{
					"id":    user.ID,
					"name":  user.Name,
					"email": user.Email,
					"role":  user.Role,
				},
				Token:     accessToken,
				ExpiresIn: tokenPair.ExpiresIn,
			},
		}
		ctx.JSON(http.StatusOK, response)
		return
	}

	// Default to cookie-based auth
	ctx.SetCookie(
		"jwt-id",
		strconv.FormatInt(tokenPair.AccessTokenID, 10),
		int(tokenPair.ExpiresIn),
		"/",
		"",
		false, // Secure
		true,  // HttpOnly
	)

	ctx.SetCookie(
		"refresh-id",
		strconv.FormatInt(tokenPair.RefreshTokenID, 10),
		int(c.jwtService.config.RefreshExpiry/time.Second),
		"/",
		"",
		false, // Secure
		true,  // HttpOnly
	)

	response := LoginResponse{
		Status:  "success",
		Message: "Login successful",
		Data: CookieResponse{
			User: gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
			ExpiresIn: tokenPair.ExpiresIn,
		},
	}
	ctx.JSON(http.StatusOK, response)
}

// Refresh generates a new access token using a refresh token
func (c *AuthController) Refresh(ctx *gin.Context) {
	// Get refresh token ID from cookie
	refreshIDStr, err := ctx.Cookie("refresh-id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token ID cookie required"})
		return
	}

	// Convert refresh token ID to int64
	refreshID, err := strconv.ParseInt(refreshIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token ID format"})
		return
	}

	// Generate new token pair
	tokenPair, err := c.jwtService.RefreshAccessToken(refreshID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Set new access token ID cookie
	ctx.SetCookie(
		"jwt-id",
		strconv.FormatInt(tokenPair.AccessTokenID, 10),
		int(tokenPair.ExpiresIn),
		"/",
		"",
		false, // Secure
		true,  // HttpOnly
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Token refreshed successfully",
		"expires_in": tokenPair.ExpiresIn,
	})
}

// Logout invalidates the current token
func (c *AuthController) Logout(ctx *gin.Context) {
	// Get token ID from context
	tokenID, exists := ctx.Get("token_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No active session"})
		return
	}

	// Delete the token from Redis
	err := DeleteToken(tokenID.(int64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate token"})
		return
	}

	// Clear cookies
	ctx.SetCookie("jwt-id", "", -1, "/", "", false, true)
	ctx.SetCookie("refresh-id", "", -1, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RegisterRoutes sets up auth endpoints
func (c *AuthController) RegisterRoutes(router gin.IRouter) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", c.Login)
		auth.POST("/refresh", c.Refresh)
		auth.POST("/logout", c.Logout)
	}
}
