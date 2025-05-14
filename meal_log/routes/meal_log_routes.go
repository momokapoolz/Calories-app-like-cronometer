package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/meal_log/controllers"
	"github.com/momokapoolz/caloriesapp/meal_log/repository"
	"github.com/momokapoolz/caloriesapp/meal_log/services"
	"gorm.io/gorm"
)

// SetupMealLogRoutes initializes meal log routes
func SetupMealLogRoutes(router *gin.RouterGroup, db *gorm.DB) {
	mealLogRepo := repository.NewMealLogRepository(db)
	mealLogService := services.NewMealLogService(mealLogRepo)
	mealLogController := controllers.NewMealLogController(mealLogService)

	mealLogRoutes := router.Group("/meal-logs")
	{
		mealLogRoutes.POST("/", mealLogController.CreateMealLog)
		mealLogRoutes.GET("/:id", mealLogController.GetMealLog)
		mealLogRoutes.GET("/user/:userId", mealLogController.GetMealLogsByUserID)
		mealLogRoutes.GET("/user/:userId/date/:date", mealLogController.GetMealLogsByUserIDAndDate)
		mealLogRoutes.GET("/user/:userId/date-range", mealLogController.GetMealLogsByUserIDAndDateRange)
		mealLogRoutes.PUT("/:id", mealLogController.UpdateMealLog)
		mealLogRoutes.DELETE("/:id", mealLogController.DeleteMealLog)
	}
} 