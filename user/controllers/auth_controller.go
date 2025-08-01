package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/dto"
	"github.com/momokapoolz/caloriesapp/user/models"
	"github.com/momokapoolz/caloriesapp/user/repository"
	"github.com/momokapoolz/caloriesapp/user/utils"
)

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
	var loginReq dto.LoginRequestDTO

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

	log.Printf("[Login] Generated token pair: AccessTokenID=%s, RefreshTokenID=%s", tokenPair.AccessTokenID, tokenPair.RefreshTokenID)

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

// Logout invalidates the current user's token
func (c *UserAuthController) Logout(ctx *gin.Context) {
	log.Printf("[Logout] Attempting logout")

	// Get token ID from context (set by auth middleware)
	tokenID, exists := ctx.Get("token_id")
	if !exists {
		log.Printf("[Logout] No token ID found in context")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "No active session",
		})
		return
	}

	log.Printf("[Logout] Token ID found: %s", tokenID.(string))

	// Delete the token from Redis
	err := auth.DeleteToken(tokenID.(string))
	if err != nil {
		log.Printf("[Logout] Failed to delete token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to invalidate token",
		})
		return
	}

	log.Printf("[Logout] Token successfully invalidated")

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logged out successfully",
	})
}

// Register creates a new user account and returns authentication tokens
func (c *UserAuthController) Register(ctx *gin.Context) {
	var registerReq dto.RegisterDTO

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&registerReq); err != nil {
		log.Printf("[Register] Invalid request format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("[Register] Attempting registration for email: %s", registerReq.Email)

	// Check if user already exists
	existingUser, err := c.userRepo.FindByEmail(registerReq.Email)
	if err == nil && existingUser != nil {
		log.Printf("[Register] Email already in use: %s", registerReq.Email)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email already in use",
		})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(registerReq.Password)
	if err != nil {
		log.Printf("[Register] Failed to hash password: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to process registration",
		})
		return
	}

	// Create new user
	user := &models.User{
		Name:          registerReq.Name,
		Email:         registerReq.Email,
		PasswordHash:  hashedPassword,
		Age:           registerReq.Age,
		Gender:        registerReq.Gender,
		Weight:        registerReq.Weight,
		Height:        registerReq.Height,
		Goal:          registerReq.Goal,
		ActivityLevel: registerReq.ActivityLevel,
		CreatedAt:     time.Now(),
		Role:          "user", // Default role
	}

	// Save user to database
	if err := c.userRepo.Create(user); err != nil {
		log.Printf("[Register] Failed to create user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create user account",
			"error":   err.Error(),
		})
		return
	}

	// Return success response with user data and tokens
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User registered successfully",
		"data": gin.H{
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
		},
	})
}

// RegisterRoutes registers the auth routes
func (c *UserAuthController) RegisterRoutes(router gin.IRouter) {
	// Auth middleware for protected routes
	authMiddleware := auth.NewAuthMiddleware()

	// Public routes
	router.POST("/login", c.Login)
	router.POST("/register", c.Register)

	// Protected routes (require authentication)
	router.POST("/logout", authMiddleware.RequireAuth(), c.Logout)
}
