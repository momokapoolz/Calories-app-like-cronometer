package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/dto"
	"github.com/momokapoolz/caloriesapp/user/models"
	"github.com/momokapoolz/caloriesapp/user/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserController handles user-related endpoints
type UserController struct {
	userRepo *repository.UserRepository
}

// NewUserController creates a new user controller
func NewUserController() *UserController {
	return &UserController{
		userRepo: repository.NewUserRepository(),
	}
}

// validateUserFromContext gets and validates the user ID from context
func (c *UserController) validateUserFromContext(ctx *gin.Context) (uint, error) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		return 0, errors.New("not authenticated")
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		return 0, errors.New("invalid user ID format")
	}

	return userID, nil
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

// toUserResponse converts a User model to UserResponseDTO
func (c *UserController) toUserResponse(user *models.User) dto.UserResponseDTO {
	return dto.UserResponseDTO{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Age:           user.Age,
		Gender:        user.Gender,
		Weight:        user.Weight,
		Height:        user.Height,
		Goal:          user.Goal,
		ActivityLevel: user.ActivityLevel,
		Role:          user.Role,
		CreatedAt:     user.CreatedAt,
	}
}

// UpdateProfile allows authenticated users to update their own profile
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	// Validate user from context
	userID, err := c.validateUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var updateReq dto.UserUpdateProfileRequestDTO

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&updateReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Fetch current user data
	user, err := c.userRepo.FindByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	// Check if email is being updated and if it already exists
	if updateReq.Email != nil && *updateReq.Email != user.Email {
		existingUser, err := c.userRepo.FindByEmail(*updateReq.Email)
		if err == nil && existingUser.ID != user.ID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Email already in use",
			})
			return
		}
	}

	// Update only provided fields
	if updateReq.Name != nil {
		user.Name = *updateReq.Name
	}
	if updateReq.Email != nil {
		user.Email = *updateReq.Email
	}
	if updateReq.Age != nil {
		user.Age = *updateReq.Age
	}
	if updateReq.Gender != nil {
		user.Gender = *updateReq.Gender
	}
	if updateReq.Weight != nil {
		user.Weight = *updateReq.Weight
	}
	if updateReq.Height != nil {
		user.Height = *updateReq.Height
	}
	if updateReq.Goal != nil {
		user.Goal = *updateReq.Goal
	}
	if updateReq.ActivityLevel != nil {
		user.ActivityLevel = *updateReq.ActivityLevel
	}

	// Save updated user to database
	if err := c.userRepo.Update(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update profile",
			"error":   err.Error(),
		})
		return
	}

	// Return updated user data
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile updated successfully",
		"data": gin.H{
			"user": c.toUserResponse(user),
		},
	})
}

// ChangePassword allows authenticated users to change their password
func (c *UserController) ChangePassword(ctx *gin.Context) {
	// Validate user from context
	userID, err := c.validateUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var changePasswordReq dto.UserChangePasswordRequestDTO

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&changePasswordReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Fetch current user data
	user, err := c.userRepo.FindByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(changePasswordReq.CurrentPassword)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Current password is incorrect",
		})
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to process password change",
		})
		return
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	if err := c.userRepo.Update(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update password",
			"error":   err.Error(),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password changed successfully",
	})
}

// DeleteAccount allows authenticated users to delete their own account
func (c *UserController) DeleteAccount(ctx *gin.Context) {
	// Validate user from context
	userID, err := c.validateUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Confirm user exists before deletion
	user, err := c.userRepo.FindByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	// Delete the user account
	if err := c.userRepo.Delete(user.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete account",
			"error":   err.Error(),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Account deleted successfully",
	})
}

// RegisterRoutes sets up user endpoints
func (c *UserController) RegisterRoutes(router gin.IRouter) {
	// Public routes

	// Protected routes - require authentication
	authMiddleware := auth.NewAuthMiddleware()
	protected := router.Group("")
	protected.Use(authMiddleware.RequireAuth())
	{
		// Profile management endpoints
		protected.GET("/profile", c.GetProfile)
		protected.PUT("/profile", c.UpdateProfile)
		protected.DELETE("/account", c.DeleteAccount)

		// Password management endpoint
		protected.PUT("/password", c.ChangePassword)
	}
}
