package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/user/repository"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents the user login data
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserAuthController handles user authentication
type UserAuthController struct {
	userRepo   *repository.UserRepository
	jwtService *auth.JWTService
}

// NewUserAuthController creates a new auth controller
func NewUserAuthController() *UserAuthController {
	return &UserAuthController{
		userRepo:   repository.NewUserRepository(),
		jwtService: auth.NewJWTService(),
	}
}

// Login authenticates a user and returns JWT tokens
func (c *UserAuthController) Login(ctx *gin.Context) {
	var loginReq LoginRequest

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Find user by email
	user, err := c.userRepo.FindByEmail(loginReq.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid credentials",
		})
		return
	}

	// Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid credentials",
		})
		return
	}

	// Generate JWT token
	tokenPair, err := c.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to generate authentication token",
		})
		return
	}

	// Return tokens to client
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login successful",
		"data": gin.H{
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
			"tokens": tokenPair,
		},
	})
}

// RegisterRoutes registers the auth routes
func (c *UserAuthController) RegisterRoutes(router *gin.Engine) {
	router.POST("/login", c.Login)
} 