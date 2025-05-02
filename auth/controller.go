package auth

import (
	"net/http"

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

	ctx.JSON(http.StatusOK, tokenPair)
}

// Refresh generates a new access token using a refresh token
func (c *AuthController) Refresh(ctx *gin.Context) {
	var request RefreshRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new token pair
	tokenPair, err := c.jwtService.RefreshAccessToken(request.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, tokenPair)
}

// RegisterRoutes sets up auth endpoints
func (c *AuthController) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", c.Login)
		auth.POST("/refresh", c.Refresh)
	}
}
