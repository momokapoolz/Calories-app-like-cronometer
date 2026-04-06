package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/momokapoolz/caloriesapp/auth"

	"github.com/momokapoolz/caloriesapp/dto"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/helpers"
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

// CreateMealLog godoc
// @Summary      Create meal log
// @Description  Create a new meal log with items for the authenticated user
// @Tags         meal_log
// @Accept       json
// @Produce      json
// @Param        meal_log  body      dto.CreateMealLogRequestDTO  true  "Meal log data"
// @Success      201  {object}  dto.CreateMealLogRequestDTO  "Meal log created successfully"
// @Failure      400  {object}  map[string]string            "Invalid request body"
// @Failure      401  {object}  map[string]string            "Unauthorized"
// @Failure      500  {object}  map[string]string            "Internal server error"
// @Security     BearerAuth
// @Router       /meal-logs/ [post]
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
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, mealLogWithItems)

}

// GetMealLog godoc
// @Summary      Get a meal log by ID
// @Description  Retrieve a specific meal log record by its ID
// @Tags         meal_log
// @Produce      json
// @Param        id  path      int  true  "Meal log ID"
// @Success      200  {object}  dto.CreateMealLogRequestDTO  "Meal log retrieved successfully"
// @Failure      400  {object}  map[string]string            "Invalid ID format"
// @Failure      404  {object}  map[string]string            "Meal log not found"
// @Failure      500  {object}  map[string]string            "Internal server error"
// @Security     BearerAuth
// @Router       /meal-logs/{id} [get]
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

// GetMealLogsByUserID godoc
// @Summary      Get meal logs by authenticated user
// @Description  Retrieve all meal logs belonging to the currently authenticated user
// @Tags         meal_log
// @Produce      json
// @Success      200  {array}   dto.CreateMealLogRequestDTO  "List of meal logs"
// @Failure      401  {object}  map[string]string            "Unauthorized"
// @Failure      500  {object}  map[string]string            "Internal server error"
// @Security     BearerAuth
// @Router       /meal-logs/user [get]
func (c *MealLogController) GetMealLogsByUserID(ctx *gin.Context) {

	claims, ok := auth.GetCurrentUser(ctx)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	mealLogs, err := c.service.GetMealLogsByUserID(uint(claims.UserID))
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal logs"})
		return
	}

	ctx.JSON(http.StatusOK, mealLogs)
}

// GetMealLogsByUserIDAndDate godoc
// @Summary      Get meal logs by date
// @Description  Retrieve all meal logs for the authenticated user on a specific date
// @Tags         meal_log
// @Produce      json
// @Param        date  path  string  true  "Date in YYYY-MM-DD format"
// @Success      200  {array}   dto.CreateMealLogRequestDTO  "List of meal logs for the date"
// @Failure      400  {object}  map[string]string            "Invalid date format"
// @Failure      401  {object}  map[string]string            "Unauthorized"
// @Failure      500  {object}  map[string]string            "Internal server error"
// @Security     BearerAuth
// @Router       /meal-logs/user/date/{date} [get]
func (c *MealLogController) GetMealLogsByUserIDAndDate(ctx *gin.Context) {
	claims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := claims.UserID

	dateStr := ctx.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	mealLogs, err := c.service.GetMealLogsByUserIDAndDate(userID, date)
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal logs for the specified date"})
		return
	}

	ctx.JSON(http.StatusOK, mealLogs)
}

// GetMealLogsByUserIDAndDateRange godoc
// @Summary      Get meal logs by date range
// @Description  Retrieve all meal logs for the authenticated user within a date range
// @Tags         meal_log
// @Produce      json
// @Param        startDate  query  string  true  "Start date in YYYY-MM-DD format"
// @Param        endDate    query  string  true  "End date in YYYY-MM-DD format"
// @Success      200  {array}   dto.CreateMealLogRequestDTO  "List of meal logs in the date range"
// @Failure      400  {object}  map[string]string            "Invalid date format or missing parameters"
// @Failure      401  {object}  map[string]string            "Unauthorized"
// @Failure      500  {object}  map[string]string            "Internal server error"
// @Security     BearerAuth
// @Router       /meal-logs/user/date-range [get]
func (c *MealLogController) GetMealLogsByUserIDAndDateRange(ctx *gin.Context) {
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
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal logs for the specified date range"})
		return
	}

	ctx.JSON(http.StatusOK, mealLogs)
}

// UpdateMealLog godoc
// @Summary      Update a meal log
// @Description  Update the meal type of an existing meal log by ID (ownership required)
// @Tags         meal_log
// @Accept       json
// @Produce      json
// @Param        id        path  int                      true  "Meal log ID"
// @Param        meal_log  body  object{meal_type=string} true  "Updated meal log data"
// @Success      200  {object}  dto.CreateMealLogRequestDTO  "Meal log updated successfully"
// @Failure      400  {object}  map[string]string            "Invalid ID or request body"
// @Failure      401  {object}  map[string]string            "Unauthorized"
// @Failure      403  {object}  map[string]string            "Forbidden — not the owner"
// @Failure      404  {object}  map[string]string            "Meal log not found"
// @Failure      500  {object}  map[string]string            "Internal server error"
// @Security     BearerAuth
// @Router       /meal-logs/{id} [put]
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
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meal log"})
		return
	}

	ctx.JSON(http.StatusOK, existingMealLog)
}

// DeleteMealLog godoc
// @Summary      Delete a meal log
// @Description  Delete a meal log by ID (ownership required)
// @Tags         meal_log
// @Produce      json
// @Param        id  path  int  true  "Meal log ID"
// @Success      200  {object}  map[string]string  "Meal log deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid ID format"
// @Failure      401  {object}  map[string]string  "Unauthorized"
// @Failure      403  {object}  map[string]string  "Forbidden — not the owner"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /meal-logs/{id} [delete]
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

	mealLog, err := c.service.GetMealLogByID(uint(id))
	if err != nil || mealLog.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this meal log"})
		return
	}

	if err := c.service.DeleteMealLog(uint(id)); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meal log"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Meal log deleted successfully"})
}
