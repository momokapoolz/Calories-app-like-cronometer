package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/user_biometrics/models"
	"github.com/momokapoolz/caloriesapp/user_biometrics/services"
)

// UserBiometricController handles HTTP requests for user biometric operations
type UserBiometricController struct {
	service *services.UserBiometricService
}

// NewUserBiometricController creates a new user biometric controller instance
func NewUserBiometricController(service *services.UserBiometricService) *UserBiometricController {
	return &UserBiometricController{service: service}
}

// CreateUserBiometric handles the creation of a new user biometric record
func (c *UserBiometricController) CreateUserBiometric(ctx *gin.Context) {
	var biometric models.UserBiometric
	if err := ctx.ShouldBindJSON(&biometric); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set created_at to current time if not provided
	if biometric.CreatedAt.IsZero() {
		biometric.CreatedAt = time.Now()
	}

	if err := c.service.CreateUserBiometric(&biometric); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user biometric"})
		return
	}

	ctx.JSON(http.StatusCreated, biometric)
}

// GetUserBiometric retrieves a user biometric by its ID
func (c *UserBiometricController) GetUserBiometric(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	biometric, err := c.service.GetUserBiometricByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User biometric not found"})
		return
	}

	ctx.JSON(http.StatusOK, biometric)
}

// GetUserBiometricsByUserID retrieves all biometrics for a specific user
func (c *UserBiometricController) GetUserBiometricsByUserID(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	biometrics, err := c.service.GetUserBiometricsByUserID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user biometrics"})
		return
	}

	ctx.JSON(http.StatusOK, biometrics)
}

// GetUserBiometricsByUserIDAndType retrieves biometrics of a specific type for a specific user
func (c *UserBiometricController) GetUserBiometricsByUserIDAndType(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	biometricType := ctx.Param("type")
	biometrics, err := c.service.GetUserBiometricsByUserIDAndType(uint(userID), biometricType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user biometrics by type"})
		return
	}

	ctx.JSON(http.StatusOK, biometrics)
}

// GetUserBiometricsByUserIDAndTypeAndDateRange retrieves biometrics of a specific type for a user within a date range
func (c *UserBiometricController) GetUserBiometricsByUserIDAndTypeAndDateRange(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	biometricType := ctx.Param("type")
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Start date and end date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Set end date to the end of the day
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	biometrics, err := c.service.GetUserBiometricsByUserIDAndTypeAndDateRange(uint(userID), biometricType, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user biometrics for the specified date range"})
		return
	}

	ctx.JSON(http.StatusOK, biometrics)
}

// GetLatestUserBiometricByUserIDAndType retrieves the most recent biometric of a specific type for a user
func (c *UserBiometricController) GetLatestUserBiometricByUserIDAndType(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	biometricType := ctx.Param("type")
	biometric, err := c.service.GetLatestUserBiometricByUserIDAndType(uint(userID), biometricType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Latest user biometric not found"})
		return
	}

	ctx.JSON(http.StatusOK, biometric)
}

// UpdateUserBiometric updates a user biometric record
func (c *UserBiometricController) UpdateUserBiometric(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var biometric models.UserBiometric
	if err := ctx.ShouldBindJSON(&biometric); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	biometric.ID = uint(id)
	if err := c.service.UpdateUserBiometric(&biometric); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user biometric"})
		return
	}

	ctx.JSON(http.StatusOK, biometric)
}

// DeleteUserBiometric removes a user biometric record
func (c *UserBiometricController) DeleteUserBiometric(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.DeleteUserBiometric(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user biometric"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User biometric deleted successfully"})
} 