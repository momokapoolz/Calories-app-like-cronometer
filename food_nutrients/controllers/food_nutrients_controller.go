package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/food_nutrients/models"
	"github.com/momokapoolz/caloriesapp/food_nutrients/services"
)

// FoodNutrientController handles HTTP requests for food nutrient operations
type FoodNutrientController struct {
	service *services.FoodNutrientService
}

// NewFoodNutrientController creates a new food nutrient controller instance
func NewFoodNutrientController(service *services.FoodNutrientService) *FoodNutrientController {
	return &FoodNutrientController{service: service}
}

// CreateFoodNutrient handles the creation of a new food nutrient record
func (c *FoodNutrientController) CreateFoodNutrient(ctx *gin.Context) {
	var foodNutrient models.FoodNutrient
	if err := ctx.ShouldBindJSON(&foodNutrient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateFoodNutrient(&foodNutrient); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create food nutrient"})
		return
	}

	ctx.JSON(http.StatusCreated, foodNutrient)
}

// GetFoodNutrient retrieves a food nutrient by its ID
func (c *FoodNutrientController) GetFoodNutrient(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	foodNutrient, err := c.service.GetFoodNutrientByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Food nutrient not found"})
		return
	}

	ctx.JSON(http.StatusOK, foodNutrient)
}

// GetAllFoodNutrients retrieves all food nutrient records
func (c *FoodNutrientController) GetAllFoodNutrients(ctx *gin.Context) {
	foodNutrients, err := c.service.GetAllFoodNutrients()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve food nutrients"})
		return
	}

	ctx.JSON(http.StatusOK, foodNutrients)
}

// GetFoodNutrientsByFoodID retrieves food nutrients by food ID
func (c *FoodNutrientController) GetFoodNutrientsByFoodID(ctx *gin.Context) {
	idStr := ctx.Param("foodId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID format"})
		return
	}

	foodNutrients, err := c.service.GetFoodNutrientsByFoodID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve food nutrients by food ID"})
		return
	}

	ctx.JSON(http.StatusOK, foodNutrients)
}

// GetFoodNutrientsByNutrientID retrieves food nutrients by nutrient ID
func (c *FoodNutrientController) GetFoodNutrientsByNutrientID(ctx *gin.Context) {
	idStr := ctx.Param("nutrientId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nutrient ID format"})
		return
	}

	foodNutrients, err := c.service.GetFoodNutrientsByNutrientID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve food nutrients by nutrient ID"})
		return
	}

	ctx.JSON(http.StatusOK, foodNutrients)
}

// UpdateFoodNutrient updates a food nutrient record
func (c *FoodNutrientController) UpdateFoodNutrient(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var foodNutrient models.FoodNutrient
	if err := ctx.ShouldBindJSON(&foodNutrient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foodNutrient.ID = uint(id)
	if err := c.service.UpdateFoodNutrient(&foodNutrient); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update food nutrient"})
		return
	}

	ctx.JSON(http.StatusOK, foodNutrient)
}

// DeleteFoodNutrient removes a food nutrient record
func (c *FoodNutrientController) DeleteFoodNutrient(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.DeleteFoodNutrient(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete food nutrient"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Food nutrient deleted successfully"})
} 