package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/food/models"
	"github.com/momokapoolz/caloriesapp/food/services"
)

// FoodController handles HTTP requests for food operations
type FoodController struct {
	service *services.FoodService
}

// NewFoodController creates a new food controller instance
func NewFoodController(service *services.FoodService) *FoodController {
	return &FoodController{service: service}
}

// CreateFood handles the creation of a new food record
func (c *FoodController) CreateFood(ctx *gin.Context) {
	var food models.Food
	if err := ctx.ShouldBindJSON(&food); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateFood(&food); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create food"})
		return
	}

	ctx.JSON(http.StatusCreated, food)
}

// GetFood retrieves a food by its ID
func (c *FoodController) GetFood(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	food, err := c.service.GetFoodByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
		return
	}

	ctx.JSON(http.StatusOK, food)
}

// GetAllFoods retrieves all food records
func (c *FoodController) GetAllFoods(ctx *gin.Context) {
	foods, err := c.service.GetAllFoods()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve foods"})
		return
	}

	ctx.JSON(http.StatusOK, foods)
}

// UpdateFood updates a food record
func (c *FoodController) UpdateFood(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var food models.Food
	if err := ctx.ShouldBindJSON(&food); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	food.ID = uint(id)
	if err := c.service.UpdateFood(&food); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update food"})
		return
	}

	ctx.JSON(http.StatusOK, food)
}

// DeleteFood removes a food record
func (c *FoodController) DeleteFood(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.DeleteFood(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete food"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Food deleted successfully"})
} 