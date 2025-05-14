package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/food_nutrients/controllers"
	"github.com/momokapoolz/caloriesapp/food_nutrients/repository"
	"github.com/momokapoolz/caloriesapp/food_nutrients/services"
	"gorm.io/gorm"
)

// SetupFoodNutrientRoutes initializes food nutrient routes
func SetupFoodNutrientRoutes(router *gin.RouterGroup, db *gorm.DB) {
	foodNutrientRepo := repository.NewFoodNutrientRepository(db)
	foodNutrientService := services.NewFoodNutrientService(foodNutrientRepo)
	foodNutrientController := controllers.NewFoodNutrientController(foodNutrientService)

	foodNutrientRoutes := router.Group("/food-nutrients")
	{
		foodNutrientRoutes.POST("/", foodNutrientController.CreateFoodNutrient)
		foodNutrientRoutes.GET("/", foodNutrientController.GetAllFoodNutrients)
		foodNutrientRoutes.GET("/:id", foodNutrientController.GetFoodNutrient)
		foodNutrientRoutes.GET("/food/:foodId", foodNutrientController.GetFoodNutrientsByFoodID)
		foodNutrientRoutes.GET("/nutrient/:nutrientId", foodNutrientController.GetFoodNutrientsByNutrientID)
		foodNutrientRoutes.PUT("/:id", foodNutrientController.UpdateFoodNutrient)
		foodNutrientRoutes.DELETE("/:id", foodNutrientController.DeleteFoodNutrient)
	}
} 