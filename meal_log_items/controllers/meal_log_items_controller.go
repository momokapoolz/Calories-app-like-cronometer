package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/dto"
	"github.com/momokapoolz/caloriesapp/helpers"
	"github.com/momokapoolz/caloriesapp/meal_log_items/models"
	"github.com/momokapoolz/caloriesapp/meal_log_items/services"
)

type MealLogItemController struct {
	service *services.MealLogItemService
}

// NewMealLogItemController creates a new meal log item controller instance
func NewMealLogItemController(service *services.MealLogItemService) *MealLogItemController {
	return &MealLogItemController{service: service}
}

// CreateMealLogItem godoc
// @Summary      Create meal log item
// @Description  Create a new meal log item record
// @Tags         meal_log_item
// @Accept       json
// @Produce      json
// @Param        meal_log_item  body      models.MealLogItem  true  "Meal log item data"
// @Success      201  {object}  models.MealLogItem  "Meal log item created successfully"
// @Failure      400  {object}  map[string]string   "Invalid request body"
// @Failure      500  {object}  map[string]string   "Internal server error"
// @Security     BearerAuth
// @Router       /meal-log-items/ [post]
func (c *MealLogItemController) CreateMealLogItem(ctx *gin.Context) {
	var item models.MealLogItem
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateMealLogItem(&item); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal log item"})
		return
	}

	ctx.JSON(http.StatusCreated, item)
}

// GetMealLogItem godoc
// @Summary      Get meal log item by ID
// @Description  Retrieve a specific meal log item record by its ID
// @Tags         meal_log_item
// @Produce      json
// @Param        id  path      int  true  "Meal log item ID"
// @Success      200  {object}  models.MealLogItem  "Meal log item retrieved successfully"
// @Failure      400  {object}  map[string]string   "Invalid ID format"
// @Failure      404  {object}  map[string]string   "Meal log item not found"
// @Failure      500  {object}  map[string]string   "Internal server error"
// @Security     BearerAuth
// @Router       /meal-log-items/{id} [get]
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

// GetMealLogItemsByMealLogID godoc
// @Summary      Get meal log items by meal log ID
// @Description  Retrieve all items belonging to a specific meal log
// @Tags         meal_log_item
// @Produce      json
// @Param        mealLogId  path      int  true  "Meal log ID"
// @Success      200  {array}   models.MealLogItem  "List of meal log items"
// @Failure      400  {object}  map[string]string   "Invalid meal log ID format"
// @Failure      500  {object}  map[string]string   "Internal server error"
// @Security     BearerAuth
// @Router       /meal-log-items/meal-log/{mealLogId} [get]
func (c *MealLogItemController) GetMealLogItemsByMealLogID(ctx *gin.Context) {
	mealLogIDStr := ctx.Param("mealLogId")
	mealLogID, err := strconv.ParseUint(mealLogIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal log ID format"})
		return
	}

	items, err := c.service.GetMealLogItemsByMealLogID(uint(mealLogID))
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal log items"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

// GetMealLogItemsByFoodID godoc
// @Summary      Get meal log items by food ID
// @Description  Retrieve all meal log items associated with a specific food
// @Tags         meal_log_item
// @Produce      json
// @Param        foodId  path      int  true  "Food ID"
// @Success      200  {array}   models.MealLogItem  "List of meal log items for the food"
// @Failure      400  {object}  map[string]string   "Invalid food ID format"
// @Failure      500  {object}  map[string]string   "Internal server error"
// @Security     BearerAuth
// @Router       /meal-log-items/food/{foodId} [get]
func (c *MealLogItemController) GetMealLogItemsByFoodID(ctx *gin.Context) {
	foodIDStr := ctx.Param("foodId")
	foodID, err := strconv.ParseUint(foodIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID format"})
		return
	}

	items, err := c.service.GetMealLogItemsByFoodID(uint(foodID))
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meal log items by food ID"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

// UpdateMealLogItem godoc
// @Summary      Update meal log item
// @Description  Update an existing meal log item by ID
// @Tags         meal_log_item
// @Accept       json
// @Produce      json
// @Param        id             path  int                true  "Meal log item ID"
// @Param        meal_log_item  body  models.MealLogItem true  "Updated meal log item data"
// @Success      200  {object}  models.MealLogItem  "Meal log item updated successfully"
// @Failure      400  {object}  map[string]string   "Invalid ID or request body"
// @Failure      500  {object}  map[string]string   "Internal server error"
// @Security     BearerAuth
// @Router       /meal-log-items/{id} [put]
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
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meal log item"})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

// DeleteMealLogItem godoc
// @Summary      Delete meal log item
// @Description  Delete a specific meal log item by ID
// @Tags         meal_log_item
// @Produce      json
// @Param        id  path  int  true  "Meal log item ID"
// @Success      200  {object}  map[string]string  "Meal log item deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid ID format"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /meal-log-items/{id} [delete]
func (c *MealLogItemController) DeleteMealLogItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.DeleteMealLogItem(uint(id)); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meal log item"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Meal log item deleted successfully"})
}

// DeleteMealLogItemsByMealLogID godoc
// @Summary      Delete all items for a meal log
// @Description  Delete all meal log items associated with a specific meal log ID
// @Tags         meal_log_item
// @Produce      json
// @Param        mealLogId  path  int  true  "Meal log ID"
// @Success      200  {object}  map[string]string  "All meal log items deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid meal log ID format"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /meal-log-items/meal-log/{mealLogId} [delete]
func (c *MealLogItemController) DeleteMealLogItemsByMealLogID(ctx *gin.Context) {
	mealLogIDStr := ctx.Param("mealLogId")
	mealLogID, err := strconv.ParseUint(mealLogIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal log ID format"})
		return
	}

	if err := c.service.DeleteMealLogItemsByMealLogID(uint(mealLogID)); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meal log items"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "All items for meal log deleted successfully"})
}

// AddItemsToMealLog godoc
// @Summary      Add items to an existing meal log
// @Description  Add one or more food items to an existing meal log (ownership required)
// @Tags         meal_log_item
// @Accept       json
// @Produce      json
// @Param        id             path  int                          true  "Meal log ID"
// @Param        meal_log_item  body  dto.AddItemsToMealLogRequestDTO  true  "Items to add"
// @Success      201  {array}   models.MealLogItem  "Items added to meal log successfully"
// @Failure      400  {object}  map[string]string   "Invalid ID or request body"
// @Failure      401  {object}  map[string]string   "Unauthorized"
// @Failure      403  {object}  map[string]string   "Forbidden — not the owner"
// @Failure      500  {object}  map[string]string   "Internal server error"
// @Security     BearerAuth
// @Router       /meal-logs/{id}/items [post]
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
			helpers.LogError(err)
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
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add items to meal log: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdItems)
}
