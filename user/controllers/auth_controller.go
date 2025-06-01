package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/user/repository"
	"github.com/momokapoolz/caloriesapp/user/utils"
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
		log.Printf("[Login] Invalid request format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("[Login] Attempting login for email: %s", loginReq.Email)

	// Find user by email
	user, err := c.userRepo.FindByEmail(loginReq.Email)
	if err != nil {
		log.Printf("[Login] User not found: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid credentials",
		})
		return
	}

	log.Printf("[Login] Found user: ID=%d, Email=%s", user.ID, user.Email)

	// Compare password with stored hash using the utility function
	err = utils.ComparePasswords(user.PasswordHash, loginReq.Password)
	if err != nil {
		log.Printf("[Login] Password comparison failed: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid credentials",
		})
		return
	}

	log.Printf("[Login] Password verified successfully")

	// Generate JWT token
	tokenPair, err := c.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		log.Printf("[Login] Failed to generate token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to generate authentication token",
		})
		return
	}

	log.Printf("[Login] Generated token pair: AccessTokenID=%d, RefreshTokenID=%d", tokenPair.AccessTokenID, tokenPair.RefreshTokenID)

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
func (c *UserAuthController) RegisterRoutes(router gin.IRouter) {
	router.POST("/login", c.Login)
}
