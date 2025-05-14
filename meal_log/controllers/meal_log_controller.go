package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/meal_log/models"
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

// CreateMealLog handles the creation of a new meal log record
func (c *MealLogController) CreateMealLog(ctx *gin.Context) {
	var mealLog models.MealLog
	if err := ctx.ShouldBindJSON(&mealLog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set created_at to current time if not provided
	if mealLog.CreatedAt.IsZero() {
		mealLog.CreatedAt = time.Now()
	}

	if err := c.service.CreateMealLog(&mealLog); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal log"})
		return
	}

	ctx.JSON(http.StatusCreated, mealLog)
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
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	mealLogs, err := c.service.GetMealLogsByUserID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal logs"})
		return
	}

	ctx.JSON(http.StatusOK, mealLogs)
}

// GetMealLogsByUserIDAndDate retrieves meal logs for a specific user on a specific date
func (c *MealLogController) GetMealLogsByUserIDAndDate(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	dateStr := ctx.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	mealLogs, err := c.service.GetMealLogsByUserIDAndDate(uint(userID), date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal logs for the specified date"})
		return
	}

	ctx.JSON(http.StatusOK, mealLogs)
}

// GetMealLogsByUserIDAndDateRange retrieves meal logs for a specific user within a date range
func (c *MealLogController) GetMealLogsByUserIDAndDateRange(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

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

	mealLogs, err := c.service.GetMealLogsByUserIDAndDateRange(uint(userID), startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal logs for the specified date range"})
		return
	}

	ctx.JSON(http.StatusOK, mealLogs)
}

// UpdateMealLog updates a meal log record
func (c *MealLogController) UpdateMealLog(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var mealLog models.MealLog
	if err := ctx.ShouldBindJSON(&mealLog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mealLog.ID = uint(id)
	if err := c.service.UpdateMealLog(&mealLog); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meal log"})
		return
	}

	ctx.JSON(http.StatusOK, mealLog)
}

// DeleteMealLog removes a meal log record
func (c *MealLogController) DeleteMealLog(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.DeleteMealLog(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meal log"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Meal log deleted successfully"})
} 