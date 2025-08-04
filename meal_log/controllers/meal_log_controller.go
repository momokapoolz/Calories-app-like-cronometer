package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/momokapoolz/caloriesapp/auth"

	"github.com/momokapoolz/caloriesapp/dto"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/meal_log/services"
)

// MealLogController handles HTTP requests for meal log operations
type MealLogController struct {
	service *services.MealLogService
}

// NewMealLogController creates a new meal log controller instance
func NewMealLogController(service *services.MealLogService) *MealLogController {
	return &MealLogController{service: service}
}

// CreateMealLog handles the creation of a new meal log record (Create meal log API using DTO)
func (c *MealLogController) CreateMealLog(ctx *gin.Context) {
	var req dto.CreateMealLogRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userClaims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	mealLogWithItems, err := c.service.CreateMealLogComprehensive(userClaims.UserID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, mealLogWithItems)

}

// GetMealLog retrieves a meal log by its ID
func (c *MealLogController) GetMealLog(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	mealLog, err := c.service.GetMealLogByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Meal log not found"})
		return
	}

	ctx.JSON(http.StatusOK, mealLog)
}

// GetMealLogsByUserID retrieves all meal logs for a specific user
func (c *MealLogController) GetMealLogsByUserID(ctx *gin.Context) {

	//update auth
	claims, ok := auth.GetCurrentUser(ctx)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	mealLogs, err := c.service.GetMealLogsByUserID(uint(claims.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal logs"})
		return
	}

	//userIDStr := ctx.Param("userId")
	//userID, err := strconv.ParseUint(userIDStr, 10, 32)
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
	//	return
	//}
	//
	//mealLogs, err := c.service.GetMealLogsByUserID(uint(userID))
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal logs"})
	//	return
	//}

	ctx.JSON(http.StatusOK, mealLogs)
}

// GetMealLogsByUserIDAndDate retrieves meal logs for a specific user on a specific date
func (c *MealLogController) GetMealLogsByUserIDAndDate(ctx *gin.Context) {
	// Lấy user từ JWT
	claims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := claims.UserID

	// Lấy và parse ngày từ path param
	dateStr := ctx.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Gọi service để lấy meal logs theo user và ngày
	mealLogs, err := c.service.GetMealLogsByUserIDAndDate(userID, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal logs for the specified date"})
		return
	}

	ctx.JSON(http.StatusOK, mealLogs)
}

// from here
// GetMealLogsByUserIDAndDateRange retrieves meal logs for a specific user within a date range
func (c *MealLogController) GetMealLogsByUserIDAndDateRange(ctx *gin.Context) {
	// Lấy user từ JWT
	claims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := claims.UserID

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

	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	mealLogs, err := c.service.GetMealLogsByUserIDAndDateRange(userID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal logs for the specified date range"})
		return
	}

	ctx.JSON(http.StatusOK, mealLogs)
}

// UpdateMealLog updates a meal log record
func (c *MealLogController) UpdateMealLog(ctx *gin.Context) {
	claims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := claims.UserID

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// First, get the existing meal log to preserve created_at
	existingMealLog, err := c.service.GetMealLogByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Meal log not found"})
		return
	}

	// Check ownership
	if existingMealLog.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this meal log"})
		return
	}

	// Create update request struct to only accept specific fields
	var updateRequest struct {
		MealType string `json:"meal_type"`
	}

	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only the meal_type, preserve other fields
	existingMealLog.MealType = updateRequest.MealType

	if err := c.service.UpdateMealLog(existingMealLog); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meal log"})
		return
	}

	ctx.JSON(http.StatusOK, existingMealLog)
}

// DeleteMealLog removes a meal log record
func (c *MealLogController) DeleteMealLog(ctx *gin.Context) {
	claims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := claims.UserID

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Optional: kiểm tra meal log này có thuộc về user không
	mealLog, err := c.service.GetMealLogByID(uint(id))
	if err != nil || mealLog.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this meal log"})
		return
	}

	if err := c.service.DeleteMealLog(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meal log"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Meal log deleted successfully"})
}
