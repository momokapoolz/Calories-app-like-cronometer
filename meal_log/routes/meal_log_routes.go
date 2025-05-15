package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/meal_log/controllers"
	mealLogRepo "github.com/momokapoolz/caloriesapp/meal_log/repository"
	"github.com/momokapoolz/caloriesapp/meal_log/services"
	mealLogItemsRepo "github.com/momokapoolz/caloriesapp/meal_log_items/repository"
	"gorm.io/gorm"
)

// SetupMealLogRoutes initializes meal log routes
func SetupMealLogRoutes(router *gin.RouterGroup, db *gorm.DB) {
	mealLogRepository := mealLogRepo.NewMealLogRepository(db)
	mealLogItemsRepository := mealLogItemsRepo.NewMealLogItemRepository(db)
	mealLogService := services.NewMealLogService(mealLogRepository, mealLogItemsRepository)
	mealLogController := controllers.NewMealLogController(mealLogService)

	authMiddleware := auth.NewAuthMiddleware()

	mealLogRoutes := router.Group("/meal-logs", authMiddleware.RequireAuth())
	{
		mealLogRoutes.POST("/", mealLogController.CreateMealLog)
		mealLogRoutes.GET("/:id", mealLogController.GetMealLog)
		mealLogRoutes.GET("/user", mealLogController.GetMealLogsByUserID)
		mealLogRoutes.GET("/user/date/:date", mealLogController.GetMealLogsByUserIDAndDate)
		mealLogRoutes.GET("/user/date-range", mealLogController.GetMealLogsByUserIDAndDateRange)
		mealLogRoutes.PUT("/:id", mealLogController.UpdateMealLog)
		mealLogRoutes.DELETE("/:id", mealLogController.DeleteMealLog)
	}
}
