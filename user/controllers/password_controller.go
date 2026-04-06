package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/dto"
	"github.com/momokapoolz/caloriesapp/helpers"
	"github.com/momokapoolz/caloriesapp/user/services"
)

type PasswordController struct {
	passwordService *services.PasswordService
}

// NewPasswordController creates a new password controller instance
func NewPasswordController(passwordService *services.PasswordService) *PasswordController {
	return &PasswordController{
		passwordService: passwordService,
	}
}

// UpdatePassword godoc
// @Summary      Update password
// @Description  Update the current authenticated user's password by verifying the current password first
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UpdatePasswordRequestDTO  true  "Current and new password"
// @Success      200  {object}  map[string]string  "Password updated successfully"
// @Failure      400  {object}  map[string]string  "Invalid request format"
// @Failure      401  {object}  map[string]string  "Unauthorized or current password incorrect"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /user/password/update [post]
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

	var req dto.UpdatePasswordRequestDTO
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
		helpers.LogError(err)
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

// AdminUpdatePassword godoc
// @Summary      Admin update user password
// @Description  Allow an administrator to update any user's password by email
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AdminUpdatePasswordRequestDTO  true  "Target email and new password"
// @Success      200  {object}  map[string]string  "Password updated successfully"
// @Failure      400  {object}  map[string]string  "Invalid request format"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /admin/user/password/update [post]
// AdminUpdatePassword allows administrators to update any user's password
func (c *PasswordController) AdminUpdatePassword(ctx *gin.Context) {
	var req dto.AdminUpdatePasswordRequestDTO
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
		helpers.LogError(err)
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
