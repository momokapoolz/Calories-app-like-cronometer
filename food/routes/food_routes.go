package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/food/controllers"
	"github.com/momokapoolz/caloriesapp/food/repository"
	"github.com/momokapoolz/caloriesapp/food/services"
	"gorm.io/gorm"
)

// SetupFoodRoutes initializes food routes
func SetupFoodRoutes(router *gin.RouterGroup, db *gorm.DB) {
	foodRepo := repository.NewFoodRepository(db)
	foodService := services.NewFoodService(foodRepo)
	foodController := controllers.NewFoodController(foodService)

	foodRoutes := router.Group("/foods")
	{
		foodRoutes.POST("/", foodController.CreateFood)
		foodRoutes.GET("/", foodController.GetAllFoods)
		foodRoutes.GET("/:id", foodController.GetFood)
		foodRoutes.PUT("/:id", foodController.UpdateFood)
		foodRoutes.DELETE("/:id", foodController.DeleteFood)
	}
} 