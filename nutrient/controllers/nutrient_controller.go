package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/nutrient/models"
	"github.com/momokapoolz/caloriesapp/nutrient/services"
)

type NutrientController struct {
	service *services.NutrientService
}

// NewNutrientController creates a new nutrient controller instance
func NewNutrientController(service *services.NutrientService) *NutrientController {
	return &NutrientController{service: service}
}

// CreateNutrient handles the creation of a new nutrient record
func (c *NutrientController) CreateNutrient(ctx *gin.Context) {
	var nutrient models.Nutrient
	if err := ctx.ShouldBindJSON(&nutrient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateNutrient(&nutrient); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create nutrient"})
		return
	}

	ctx.JSON(http.StatusCreated, nutrient)
}

// GetNutrient retrieves a nutrient by its ID
func (c *NutrientController) GetNutrient(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	nutrient, err := c.service.GetNutrientByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Nutrient not found"})
		return
	}

	ctx.JSON(http.StatusOK, nutrient)
}

// GetAllNutrients retrieves all nutrient records
func (c *NutrientController) GetAllNutrients(ctx *gin.Context) {
	nutrients, err := c.service.GetAllNutrients()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve nutrients"})
		return
	}

	ctx.JSON(http.StatusOK, nutrients)
}

// GetNutrientsByCategory retrieves nutrients by category
func (c *NutrientController) GetNutrientsByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	nutrients, err := c.service.GetNutrientsByCategory(category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve nutrients by category"})
		return
	}

	ctx.JSON(http.StatusOK, nutrients)
}

// UpdateNutrient updates a nutrient record
func (c *NutrientController) UpdateNutrient(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var nutrient models.Nutrient
	if err := ctx.ShouldBindJSON(&nutrient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nutrient.ID = uint(id)
	if err := c.service.UpdateNutrient(&nutrient); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update nutrient"})
		return
	}

	ctx.JSON(http.StatusOK, nutrient)
}

// DeleteNutrient removes a nutrient record
func (c *NutrientController) DeleteNutrient(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.DeleteNutrient(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete nutrient"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Nutrient deleted successfully"})
}

// GetUserNutritionByDate returns nutrition calculation for a specific date
func (c *NutrientController) GetUserNutritionByDate(ctx *gin.Context) {
	// Get authenticated user
	userClaims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse date from path parameter
	dateStr := ctx.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Get nutrition summary
	summary, err := c.service.CalculateUserNutritionByDate(userClaims.UserID, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate nutrition"})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// GetUserNutritionByDateRange returns nutrition calculation for a date range
func (c *NutrientController) GetUserNutritionByDateRange(ctx *gin.Context) {
	// Get authenticated user
	userClaims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse query parameters
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "startDate and endDate are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Use YYYY-MM-DD"})
		return
	}

	// Extend end date to include the full day
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	// Get nutrition summary
	summary, err := c.service.CalculateUserNutritionByDateRange(userClaims.UserID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate nutrition"})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// GetUserCurrentNutrition returns nutrition calculation for today
func (c *NutrientController) GetUserCurrentNutrition(ctx *gin.Context) {
	// Get authenticated user
	userClaims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Use today's date
	today := time.Now().Truncate(24 * time.Hour)

	// Get nutrition summary
	summary, err := c.service.CalculateUserNutritionByDate(userClaims.UserID, today)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate nutrition"})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// GetMealNutrition returns nutrition calculation for a specific meal
func (c *NutrientController) GetMealNutrition(ctx *gin.Context) {
	// Get authenticated user
	userClaims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse meal log ID from path parameter
	mealLogIDStr := ctx.Param("mealLogId")
	mealLogID, err := strconv.ParseUint(mealLogIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal log ID format"})
		return
	}

	// Calculate nutrition for the specific meal
	mealNutrition, err := c.service.CalculateMealNutrition(uint(mealLogID), userClaims.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate meal nutrition"})
		return
	}

	ctx.JSON(http.StatusOK, mealNutrition)
}
