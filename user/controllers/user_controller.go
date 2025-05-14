package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/user/models"
	"github.com/momokapoolz/caloriesapp/user/repository"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents the user registration data
type RegisterRequest struct {
	Name          string  `json:"name" binding:"required"`
	Email         string  `json:"email" binding:"required,email"`
	Password      string  `json:"password" binding:"required,min=6"`
	Age           int64   `json:"age" binding:"required"`
	Gender        string  `json:"gender" binding:"required"`
	Weight        float64 `json:"weight" binding:"required"`
	Height        float64 `json:"height" binding:"required"`
	Goal          string  `json:"goal" binding:"required"`
	ActivityLevel string  `json:"activity_level" binding:"required"`
}

// UserController handles user-related endpoints
type UserController struct {
	userRepo   *repository.UserRepository
	jwtService *auth.JWTService
}

// NewUserController creates a new user controller
func NewUserController() *UserController {
	return &UserController{
		userRepo:   repository.NewUserRepository(),
		jwtService: auth.NewJWTService(),
	}
}

// Register creates a new user account
func (c *UserController) Register(ctx *gin.Context) {
	var registerReq RegisterRequest

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&registerReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Check if email already exists
	_, err := c.userRepo.FindByEmail(registerReq.Email)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email already in use",
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to process registration",
		})
		return
	}

	// Create new user
	newUser := models.User{
		Name:          registerReq.Name,
		Email:         registerReq.Email,
		PasswordHash:  string(hashedPassword),
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
	if err := c.userRepo.Create(&newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create user account",
			"error":   err.Error(),
		})
		return
	}

	// Generate JWT token for the new user
	tokenPair, err := c.jwtService.GenerateTokenPair(newUser.ID, newUser.Email, newUser.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Account created but failed to generate authentication token",
		})
		return
	}

	// Return user data and tokens
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User registered successfully",
		"data": gin.H{
			"user": gin.H{
				"id":    newUser.ID,
				"name":  newUser.Name,
				"email": newUser.Email,
				"role":  newUser.Role,
			},
			"tokens": tokenPair,
		},
	})
}

// GetProfile retrieves the authenticated user's profile
func (c *UserController) GetProfile(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Not authenticated",
		})
		return
	}

	// Fetch user profile
	user, err := c.userRepo.FindByID(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	// Return user profile data
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": gin.H{
				"id":             user.ID,
				"name":           user.Name,
				"email":          user.Email,
				"age":            user.Age,
				"gender":         user.Gender,
				"weight":         user.Weight,
				"height":         user.Height,
				"goal":           user.Goal,
				"activity_level": user.ActivityLevel,
				"role":           user.Role,
				"created_at":     user.CreatedAt,
			},
		},
	})
}

// RegisterRoutes sets up user endpoints
func (c *UserController) RegisterRoutes(router gin.IRouter) {
	router.POST("/register", c.Register)

	// Protected routes
	authMiddleware := auth.NewAuthMiddleware()
	protected := router.Group("/api")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.GET("/profile", c.GetProfile)
	}
}
