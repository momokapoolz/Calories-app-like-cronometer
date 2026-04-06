package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/helpers"
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

// CreateNutrient godoc
// @Summary      Create nutrient
// @Description  Create a new nutrient record
// @Tags         nutrient
// @Accept       json
// @Produce      json
// @Param        nutrient  body      models.Nutrient  true  "Nutrient data"
// @Success      201  {object}  models.Nutrient    "Nutrient created successfully"
// @Failure      400  {object}  map[string]string  "Invalid request body"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /nutrients/ [post]
func (c *NutrientController) CreateNutrient(ctx *gin.Context) {
	var nutrient models.Nutrient
	if err := ctx.ShouldBindJSON(&nutrient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateNutrient(&nutrient); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create nutrient"})
		return
	}

	ctx.JSON(http.StatusCreated, nutrient)
}

// GetNutrient godoc
// @Summary      Get a nutrient by ID
// @Description  Retrieve a specific nutrient record by its ID
// @Tags         nutrient
// @Produce      json
// @Param        id  path      int  true  "Nutrient ID"
// @Success      200  {object}  models.Nutrient    "Nutrient retrieved successfully"
// @Failure      400  {object}  map[string]string  "Invalid ID format"
// @Failure      404  {object}  map[string]string  "Nutrient not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /nutrients/{id} [get]
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

// GetAllNutrients godoc
// @Summary      Get all nutrients
// @Description  Retrieve all nutrient records
// @Tags         nutrient
// @Produce      json
// @Success      200  {array}   models.Nutrient    "List of nutrients"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /nutrients/ [get]
func (c *NutrientController) GetAllNutrients(ctx *gin.Context) {
	nutrients, err := c.service.GetAllNutrients()
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve nutrients"})
		return
	}

	ctx.JSON(http.StatusOK, nutrients)
}

// GetNutrientsByCategory godoc
// @Summary      Get nutrients by category
// @Description  Retrieve all nutrients belonging to a specific category
// @Tags         nutrient
// @Produce      json
// @Param        category  path      string  true  "Nutrient category"
// @Success      200  {array}   models.Nutrient    "List of nutrients in the category"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /nutrients/category/{category} [get]
func (c *NutrientController) GetNutrientsByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	nutrients, err := c.service.GetNutrientsByCategory(category)
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve nutrients by category"})
		return
	}

	ctx.JSON(http.StatusOK, nutrients)
}

// UpdateNutrient godoc
// @Summary      Update nutrient
// @Description  Update an existing nutrient record by ID
// @Tags         nutrient
// @Accept       json
// @Produce      json
// @Param        id        path  int             true  "Nutrient ID"
// @Param        nutrient  body  models.Nutrient true  "Updated nutrient data"
// @Success      200  {object}  models.Nutrient    "Nutrient updated successfully"
// @Failure      400  {object}  map[string]string  "Invalid ID or request body"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /nutrients/{id} [put]
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
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update nutrient"})
		return
	}

	ctx.JSON(http.StatusOK, nutrient)
}

// DeleteNutrient godoc
// @Summary      Delete nutrient
// @Description  Delete a nutrient record by ID
// @Tags         nutrient
// @Produce      json
// @Param        id  path  int  true  "Nutrient ID"
// @Success      200  {object}  map[string]string  "Nutrient deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid ID format"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /nutrients/{id} [delete]
func (c *NutrientController) DeleteNutrient(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.DeleteNutrient(uint(id)); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete nutrient"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Nutrient deleted successfully"})
}

// GetUserNutritionByDate godoc
// @Summary      Get user nutrition by date
// @Description  Calculate and retrieve the authenticated user's total nutrition intake for a specific date
// @Tags         nutrient
// @Produce      json
// @Param        date  path  string  true  "Date in YYYY-MM-DD format"
// @Success      200  {object}  map[string]interface{}  "Nutrition summary for the date"
// @Failure      400  {object}  map[string]string       "Invalid date format"
// @Failure      401  {object}  map[string]string       "Unauthorized"
// @Failure      500  {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /nutrition/date/{date} [get]
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
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate nutrition"})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// GetUserNutritionByDateRange godoc
// @Summary      Get user nutrition by date range
// @Description  Calculate and retrieve the authenticated user's total nutrition intake within a date range
// @Tags         nutrient
// @Produce      json
// @Param        startDate  query  string  true  "Start date in YYYY-MM-DD format"
// @Param        endDate    query  string  true  "End date in YYYY-MM-DD format"
// @Success      200  {object}  map[string]interface{}  "Nutrition summary for the date range"
// @Failure      400  {object}  map[string]string       "Invalid date format or missing parameters"
// @Failure      401  {object}  map[string]string       "Unauthorized"
// @Failure      500  {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /nutrition/range [get]
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
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate nutrition"})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// GetUserCurrentNutrition godoc
// @Summary      Get today's nutrition for authenticated user
// @Description  Calculate and retrieve the authenticated user's total nutrition intake for today
// @Tags         nutrient
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Today's nutrition summary"
// @Failure      401  {object}  map[string]string       "Unauthorized"
// @Failure      500  {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /nutrition/today [get]
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
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate nutrition"})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// GetMealNutrition godoc
// @Summary      Get nutrition for a specific meal log
// @Description  Calculate and retrieve the nutrition breakdown for a specific meal log (ownership required)
// @Tags         nutrient
// @Produce      json
// @Param        mealLogId  path  int  true  "Meal log ID"
// @Success      200  {object}  map[string]interface{}  "Nutrition summary for the meal"
// @Failure      400  {object}  map[string]string       "Invalid meal log ID format"
// @Failure      401  {object}  map[string]string       "Unauthorized"
// @Failure      500  {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /nutrition/meal/{mealLogId} [get]
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
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate meal nutrition"})
		return
	}

	ctx.JSON(http.StatusOK, mealNutrition)
}
