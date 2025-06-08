package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
