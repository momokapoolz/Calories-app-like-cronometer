package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	foodRepo "github.com/momokapoolz/caloriesapp/food/repository"
	"github.com/momokapoolz/caloriesapp/meal_log_items/controllers"
	"github.com/momokapoolz/caloriesapp/meal_log_items/repository"
	"github.com/momokapoolz/caloriesapp/meal_log_items/services"
	"gorm.io/gorm"
)

// SetupMealLogItemRoutes initializes meal log item routes
func SetupMealLogItemRoutes(router *gin.RouterGroup, db *gorm.DB) {
	mealLogItemRepo := repository.NewMealLogItemRepository(db)
	foodRepository := foodRepo.NewFoodRepository(db)
	mealLogItemService := services.NewMealLogItemService(mealLogItemRepo, foodRepository)
	mealLogItemController := controllers.NewMealLogItemController(mealLogItemService)

	authMiddleware := auth.NewAuthMiddleware()

	mealLogItemRoutes := router.Group("/meal-log-items", authMiddleware.RequireAuth())
	{
		mealLogItemRoutes.POST("/", mealLogItemController.CreateMealLogItem)
		mealLogItemRoutes.GET("/:id", mealLogItemController.GetMealLogItem)
		mealLogItemRoutes.GET("/meal-log/:mealLogId", mealLogItemController.GetMealLogItemsByMealLogID)
		mealLogItemRoutes.GET("/food/:foodId", mealLogItemController.GetMealLogItemsByFoodID)
		mealLogItemRoutes.PUT("/:id", mealLogItemController.UpdateMealLogItem)
		mealLogItemRoutes.DELETE("/:id", mealLogItemController.DeleteMealLogItem)
		mealLogItemRoutes.DELETE("/meal-log/:mealLogId", mealLogItemController.DeleteMealLogItemsByMealLogID)
	}

	//Add route for adding items to a meal log with authentication middleware
	mealLogRoutes := router.Group("/meal-logs", authMiddleware.RequireAuth())
	{
		mealLogRoutes.POST("/:id/items", mealLogItemController.AddItemsToMealLog)
	}
}
