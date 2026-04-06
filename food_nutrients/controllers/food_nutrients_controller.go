package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/food_nutrients/models"
	"github.com/momokapoolz/caloriesapp/food_nutrients/services"
	"github.com/momokapoolz/caloriesapp/helpers"
)

type FoodNutrientController struct {
	service *services.FoodNutrientService
}

// NewFoodNutrientController creates a new food nutrient controller instance
func NewFoodNutrientController(service *services.FoodNutrientService) *FoodNutrientController {
	return &FoodNutrientController{service: service}
}

// CreateFoodNutrient godoc
// @Summary      Create food nutrient
// @Description  Create a new food nutrient record
// @Tags         food_nutrient
// @Accept       json
// @Produce      json
// @Param        food_nutrient  body      models.FoodNutrient  true  "Food nutrient data"
// @Success      201  {object}  models.FoodNutrient       "Food nutrient created successfully"
// @Failure      400  {object}  map[string]string         "Invalid request body"
// @Failure      500  {object}  map[string]string         "Internal server error"
// @Security     BearerAuth
// @Router       /food-nutrients/ [post]
func (c *FoodNutrientController) CreateFoodNutrient(ctx *gin.Context) {
	var foodNutrient models.FoodNutrient
	if err := ctx.ShouldBindJSON(&foodNutrient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateFoodNutrient(&foodNutrient); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create food nutrient"})
		return
	}

	ctx.JSON(http.StatusCreated, foodNutrient)
}

// GetFoodNutrient godoc
// @Summary      Get food nutrient by ID
// @Description  Retrieve a food nutrient record by its ID
// @Tags         food_nutrient
// @Produce      json
// @Param        id  path      int  true  "Food nutrient ID"
// @Success      200  {object}  models.FoodNutrient       "Food nutrient retrieved successfully"
// @Failure      400  {object}  map[string]string         "Invalid ID format"
// @Failure      404  {object}  map[string]string         "Food nutrient not found"
// @Failure      500  {object}  map[string]string         "Internal server error"
// @Security     BearerAuth
// @Router       /food-nutrients/{id} [get]
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

// GetAllFoodNutrients godoc
// @Summary      Get all food nutrients
// @Description  Retrieve all food nutrient records
// @Tags         food_nutrient
// @Produce      json
// @Success      200  {array}   models.FoodNutrient       "List of food nutrients"
// @Failure      500  {object}  map[string]string         "Internal server error"
// @Security     BearerAuth
// @Router       /food-nutrients/ [get]
func (c *FoodNutrientController) GetAllFoodNutrients(ctx *gin.Context) {
	foodNutrients, err := c.service.GetAllFoodNutrients()
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve food nutrients"})
		return
	}

	ctx.JSON(http.StatusOK, foodNutrients)
}

// GetFoodNutrientsByFoodID godoc
// @Summary      Get food nutrients by food ID
// @Description  Retrieve all nutrient records associated with a specific food
// @Tags         food_nutrient
// @Produce      json
// @Param        foodId  path      int  true  "Food ID"
// @Success      200  {array}   models.FoodNutrient       "List of food nutrients for the food"
// @Failure      400  {object}  map[string]string         "Invalid food ID format"
// @Failure      500  {object}  map[string]string         "Internal server error"
// @Security     BearerAuth
// @Router       /food-nutrients/food/{foodId} [get]
func (c *FoodNutrientController) GetFoodNutrientsByFoodID(ctx *gin.Context) {
	idStr := ctx.Param("foodId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID format"})
		return
	}

	foodNutrients, err := c.service.GetFoodNutrientsByFoodID(uint(id))
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve food nutrients by food ID"})
		return
	}

	ctx.JSON(http.StatusOK, foodNutrients)
}

// GetFoodNutrientsByNutrientID godoc
// @Summary      Get food nutrients by nutrient ID
// @Description  Retrieve all food nutrient records associated with a specific nutrient
// @Tags         food_nutrient
// @Produce      json
// @Param        nutrientId  path      int  true  "Nutrient ID"
// @Success      200  {array}   models.FoodNutrient       "List of food nutrients for the nutrient"
// @Failure      400  {object}  map[string]string         "Invalid nutrient ID format"
// @Failure      500  {object}  map[string]string         "Internal server error"
// @Security     BearerAuth
// @Router       /food-nutrients/nutrient/{nutrientId} [get]
func (c *FoodNutrientController) GetFoodNutrientsByNutrientID(ctx *gin.Context) {
	idStr := ctx.Param("nutrientId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nutrient ID format"})
		return
	}

	foodNutrients, err := c.service.GetFoodNutrientsByNutrientID(uint(id))
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve food nutrients by nutrient ID"})
		return
	}

	ctx.JSON(http.StatusOK, foodNutrients)
}

// UpdateFoodNutrient godoc
// @Summary      Update food nutrient
// @Description  Update an existing food nutrient record by ID
// @Tags         food_nutrient
// @Accept       json
// @Produce      json
// @Param        id             path      int                  true  "Food nutrient ID"
// @Param        food_nutrient  body      models.FoodNutrient  true  "Updated food nutrient data"
// @Success      200  {object}  models.FoodNutrient       "Food nutrient updated successfully"
// @Failure      400  {object}  map[string]string         "Invalid ID or request body"
// @Failure      500  {object}  map[string]string         "Internal server error"
// @Security     BearerAuth
// @Router       /food-nutrients/{id} [put]
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
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update food nutrient"})
		return
	}

	ctx.JSON(http.StatusOK, foodNutrient)
}

// DeleteFoodNutrient godoc
// @Summary      Delete food nutrient
// @Description  Delete a food nutrient record by ID
// @Tags         food_nutrient
// @Produce      json
// @Param        id  path      int  true  "Food nutrient ID"
// @Success      200  {object}  map[string]string  "Food nutrient deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid ID format"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /food-nutrients/{id} [delete]
func (c *FoodNutrientController) DeleteFoodNutrient(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.DeleteFoodNutrient(uint(id)); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete food nutrient"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Food nutrient deleted successfully"})
}
