package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/user/services"
)

type PasswordController struct {
	passwordService *services.PasswordService
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type AdminUpdatePasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func NewPasswordController(passwordService *services.PasswordService) *PasswordController {
	return &PasswordController{
		passwordService: passwordService,
	}
}

// UpdatePassword handles password update requests from users
func (c *PasswordController) UpdatePassword(ctx *gin.Context) {
	// Get current user from context (set by auth middleware)
	userClaims, exists := auth.GetCurrentUser(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Not authenticated",
		})
		return
	}

	var req UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Validate current password
	if err := c.passwordService.ValidateCurrentPassword(userClaims.Email, req.CurrentPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Current password is incorrect",
		})
		return
	}

	// Update the password
	err := c.passwordService.UpdatePassword(userClaims.Email, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update password",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password updated successfully",
	})
}

// AdminUpdatePassword allows administrators to update any user's password
func (c *PasswordController) AdminUpdatePassword(ctx *gin.Context) {
	var req AdminUpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Update the password
	err := c.passwordService.UpdatePassword(req.Email, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update password",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password updated successfully",
	})
}
