package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/food/models"
	"github.com/momokapoolz/caloriesapp/food/services"
	"github.com/momokapoolz/caloriesapp/helpers"
)

type FoodController struct {
	service *services.FoodService
}

// NewFoodController creates a new food controller instance
func NewFoodController(service *services.FoodService) *FoodController {
	return &FoodController{service: service}
}

// CreateFood godoc
// @Summary      Create new food
// @Description  Create a new food record
// @Tags         food
// @Accept       json
// @Produce      json
// @Param        food  body      models.Food       true  "Food data"
// @Success      201   {object}  models.Food       "Food created successfully"
// @Failure      400   {object}  map[string]string "Invalid request body"
// @Failure      500   {object}  map[string]string "Internal server error"
// @Security     BearerAuth
// @Router       /foods/ [post]
func (c *FoodController) CreateFood(ctx *gin.Context) {
	var food models.Food
	if err := ctx.ShouldBindJSON(&food); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateFood(&food); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create food"})
		return
	}

	ctx.JSON(http.StatusCreated, food)
}

// GetFood godoc
// @Summary      Get a specific food
// @Description  Retrieve a food record by ID
// @Tags         food
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "Food ID"
// @Success      200   {object}  models.Food       "Food retrieved successfully"
// @Failure      400   {object}  map[string]string "Invalid ID format"
// @Failure      404   {object}  map[string]string "Food not found"
// @Failure      500   {object}  map[string]string "Internal server error"
// @Security     BearerAuth
// @Router       /foods/{id} [get]
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

// GetAllFoods godoc
// @Summary      Get all foods
// @Description  Retrieve all food records
// @Tags         food
// @Produce      json
// @Success      200  {array}   models.Food       "List of foods"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Security     BearerAuth
// @Router       /foods/ [get]
func (c *FoodController) GetAllFoods(ctx *gin.Context) {
	foods, err := c.service.GetAllFoods()
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve foods"})
		return
	}

	ctx.JSON(http.StatusOK, foods)
}

// UpdateFood godoc
// @Summary      Update food
// @Description  Update an existing food record by ID
// @Tags         food
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "Food ID"
// @Param        food  body      models.Food       true  "Updated food data"
// @Success      200   {object}  models.Food       "Food updated successfully"
// @Failure      400   {object}  map[string]string "Invalid ID or request body"
// @Failure      500   {object}  map[string]string "Internal server error"
// @Security     BearerAuth
// @Router       /foods/{id} [put]
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
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update food"})
		return
	}

	ctx.JSON(http.StatusOK, food)
}

// DeleteFood godoc
// @Summary      Delete food
// @Description  Delete a food record by ID
// @Tags         food
// @Produce      json
// @Param        id   path      int               true  "Food ID"
// @Success      200  {object}  map[string]string "Food deleted successfully"
// @Failure      400  {object}  map[string]string "Invalid ID format"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Security     BearerAuth
// @Router       /foods/{id} [delete]
func (c *FoodController) DeleteFood(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.DeleteFood(uint(id)); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete food"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Food deleted successfully"})
}
