package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/dto"
	"github.com/momokapoolz/caloriesapp/meal_log_items/models"
	"github.com/momokapoolz/caloriesapp/meal_log_items/services"
)

// MealLogItemController handles HTTP requests for meal log item operations
type MealLogItemController struct {
	service *services.MealLogItemService
}

// NewMealLogItemController creates a new meal log item controller instance
func NewMealLogItemController(service *services.MealLogItemService) *MealLogItemController {
	return &MealLogItemController{service: service}
}

// CreateMealLogItem handles the creation of a new meal log item
func (c *MealLogItemController) CreateMealLogItem(ctx *gin.Context) {
	var item models.MealLogItem
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateMealLogItem(&item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal log item"})
		return
	}

	ctx.JSON(http.StatusCreated, item)
}

// GetMealLogItem retrieves a meal log item by its ID
func (c *MealLogItemController) GetMealLogItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	item, err := c.service.GetMealLogItemByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Meal log item not found"})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

// GetMealLogItemsByMealLogID retrieves all items for a specific meal log
func (c *MealLogItemController) GetMealLogItemsByMealLogID(ctx *gin.Context) {
	mealLogIDStr := ctx.Param("mealLogId")
	mealLogID, err := strconv.ParseUint(mealLogIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal log ID format"})
		return
	}

	items, err := c.service.GetMealLogItemsByMealLogID(uint(mealLogID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal log items"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

// GetMealLogItemsByFoodID retrieves all meal log items for a specific food
func (c *MealLogItemController) GetMealLogItemsByFoodID(ctx *gin.Context) {
	foodIDStr := ctx.Param("foodId")
	foodID, err := strconv.ParseUint(foodIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID format"})
		return
	}

	items, err := c.service.GetMealLogItemsByFoodID(uint(foodID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal log items by food ID"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

// UpdateMealLogItem updates a meal log item
func (c *MealLogItemController) UpdateMealLogItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var item models.MealLogItem
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.ID = uint(id)
	if err := c.service.UpdateMealLogItem(&item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meal log item"})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

// DeleteMealLogItem removes a meal log item
func (c *MealLogItemController) DeleteMealLogItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.DeleteMealLogItem(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meal log item"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Meal log item deleted successfully"})
}

// DeleteMealLogItemsByMealLogID removes all items for a specific meal log
func (c *MealLogItemController) DeleteMealLogItemsByMealLogID(ctx *gin.Context) {
	mealLogIDStr := ctx.Param("mealLogId")
	mealLogID, err := strconv.ParseUint(mealLogIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal log ID format"})
		return
	}

	if err := c.service.DeleteMealLogItemsByMealLogID(uint(mealLogID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meal log items"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "All items for meal log deleted successfully"})
}

// AddItemsToMealLog adds multiple food items to an existing meal log
func (c *MealLogItemController) AddItemsToMealLog(ctx *gin.Context) {
	// Get authenticated user from context
	userClaims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse meal log ID from path parameter
	mealLogIDStr := ctx.Param("id")
	mealLogID, err := strconv.ParseUint(mealLogIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal log ID format"})
		return
	}

	// Verify the meal log belongs to the authenticated user
	if err := c.service.VerifyMealLogOwnership(uint(mealLogID), userClaims.UserID); err != nil {
		if errors.Is(err, services.ErrUnauthorizedAccess) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to modify this meal log"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify meal log ownership: " + err.Error()})
		}
		return
	}

	// Parse request body
	var requestDTO dto.AddItemsToMealLogRequestDTO
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if len(requestDTO.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No items provided"})
		return
	}

	// Convert DTO to model objects
	var mealLogItems []models.MealLogItem
	for _, item := range requestDTO.Items {
		mealLogItems = append(mealLogItems, models.MealLogItem{
			MealLogID:     uint(mealLogID),
			FoodID:        item.FoodID,
			Quantity:      item.Quantity,
			QuantityGrams: item.QuantityGrams,
		})
	}

	// Call service to add items
	createdItems, err := c.service.AddItemsToMealLog(uint(mealLogID), mealLogItems)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add items to meal log: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdItems)
}
