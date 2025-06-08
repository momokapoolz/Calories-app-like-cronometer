package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/nutrient/controllers"
	"github.com/momokapoolz/caloriesapp/nutrient/repository"
	"github.com/momokapoolz/caloriesapp/nutrient/services"
	"gorm.io/gorm"
)

func SetupNutrientRoutes(router *gin.RouterGroup, db *gorm.DB) {
	nutrientRepo := repository.NewNutrientRepository(db)
	nutrientService := services.NewNutrientService(nutrientRepo)
	nutrientController := controllers.NewNutrientController(nutrientService)

	nutrientRoutes := router.Group("/nutrients")
	{
		nutrientRoutes.POST("/", nutrientController.CreateNutrient)
		nutrientRoutes.GET("/", nutrientController.GetAllNutrients)
		nutrientRoutes.GET("/:id", nutrientController.GetNutrient)
		nutrientRoutes.GET("/category/:category", nutrientController.GetNutrientsByCategory)
		nutrientRoutes.PUT("/:id", nutrientController.UpdateNutrient)
		nutrientRoutes.DELETE("/:id", nutrientController.DeleteNutrient)
	}
}
