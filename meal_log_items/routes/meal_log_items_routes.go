package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/meal_log_items/controllers"
	"github.com/momokapoolz/caloriesapp/meal_log_items/repository"
	"github.com/momokapoolz/caloriesapp/meal_log_items/services"
	"gorm.io/gorm"
)

// SetupMealLogItemRoutes initializes meal log item routes
func SetupMealLogItemRoutes(router *gin.RouterGroup, db *gorm.DB) {
	mealLogItemRepo := repository.NewMealLogItemRepository(db)
	mealLogItemService := services.NewMealLogItemService(mealLogItemRepo)
	mealLogItemController := controllers.NewMealLogItemController(mealLogItemService)

	mealLogItemRoutes := router.Group("/meal-log-items")
	{
		mealLogItemRoutes.POST("/", mealLogItemController.CreateMealLogItem)
		mealLogItemRoutes.GET("/:id", mealLogItemController.GetMealLogItem)
		mealLogItemRoutes.GET("/meal-log/:mealLogId", mealLogItemController.GetMealLogItemsByMealLogID)
		mealLogItemRoutes.GET("/food/:foodId", mealLogItemController.GetMealLogItemsByFoodID)
		mealLogItemRoutes.PUT("/:id", mealLogItemController.UpdateMealLogItem)
		mealLogItemRoutes.DELETE("/:id", mealLogItemController.DeleteMealLogItem)
		mealLogItemRoutes.DELETE("/meal-log/:mealLogId", mealLogItemController.DeleteMealLogItemsByMealLogID)
	}
} 