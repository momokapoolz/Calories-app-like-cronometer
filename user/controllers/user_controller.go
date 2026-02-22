package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/dto"
	"github.com/momokapoolz/caloriesapp/user/models"
	"github.com/momokapoolz/caloriesapp/user/repository"
)

// UserController handles user profile endpoints
type UserController struct {
	userRepo *repository.UserRepository
}

// NewUserController creates a new UserController
func NewUserController() *UserController {
	return &UserController{
		userRepo: repository.NewUserRepository(),
	}
}

// validateUserFromContext extracts and type-asserts the user_id set by RequireAuth
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

// toUserResponse maps a User model to the UserResponseDTO used in all profile responses
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

// GetProfile retrieves the authenticated user's full profile
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, err := c.validateUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	user, err := c.userRepo.FindByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   gin.H{"user": c.toUserResponse(user)},
	})
}

// UpdateProfile allows an authenticated user to partially update their profile.
// Only fields present in the request body are modified.
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, err := c.validateUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var req dto.UserUpdateProfileRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	user, err := c.userRepo.FindByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	// If email is changing, ensure it is not already taken by another account
	if req.Email != nil && *req.Email != user.Email {
		existing, err := c.userRepo.FindByEmail(*req.Email)
		if err == nil && existing.ID != user.ID {
			ctx.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Email already in use",
			})
			return
		}
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Age != nil {
		user.Age = *req.Age
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.Weight != nil {
		user.Weight = *req.Weight
	}
	if req.Height != nil {
		user.Height = *req.Height
	}
	if req.Goal != nil {
		user.Goal = *req.Goal
	}
	if req.ActivityLevel != nil {
		user.ActivityLevel = *req.ActivityLevel
	}

	if err := c.userRepo.Update(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update profile",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile updated successfully",
		"data":    gin.H{"user": c.toUserResponse(user)},
	})
}

// DeleteAccount permanently removes the authenticated user's account
func (c *UserController) DeleteAccount(ctx *gin.Context) {
	userID, err := c.validateUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	user, err := c.userRepo.FindByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	if err := c.userRepo.Delete(user.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete account",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Account deleted successfully",
	})
}

// RegisterRoutes registers user profile endpoints on the provided router group.
// All routes are protected by the provided authMiddleware instance.
func (c *UserController) RegisterRoutes(router gin.IRouter, authMiddleware *auth.AuthMiddleware) {
	protected := router.Group("")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.GET("/profile", c.GetProfile)
		protected.PUT("/profile", c.UpdateProfile)
		protected.DELETE("/account", c.DeleteAccount)
	}
}
