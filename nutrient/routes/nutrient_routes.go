package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	foodRepo "github.com/momokapoolz/caloriesapp/food/repository"
	foodNutrientsRepo "github.com/momokapoolz/caloriesapp/food_nutrients/repository"
	mealLogRepo "github.com/momokapoolz/caloriesapp/meal_log/repository"
	mealLogItemsRepo "github.com/momokapoolz/caloriesapp/meal_log_items/repository"
	"github.com/momokapoolz/caloriesapp/nutrient/controllers"
	"github.com/momokapoolz/caloriesapp/nutrient/repository"
	"github.com/momokapoolz/caloriesapp/nutrient/services"
	"gorm.io/gorm"
)

func SetupNutrientRoutes(router *gin.RouterGroup, db *gorm.DB) {
	// Initialize all required repositories
	nutrientRepo := repository.NewNutrientRepository(db)
	mealLogRepository := mealLogRepo.NewMealLogRepository(db)
	mealLogItemsRepository := mealLogItemsRepo.NewMealLogItemRepository(db)
	foodRepository := foodRepo.NewFoodRepository(db)
	foodNutrientsRepository := foodNutrientsRepo.NewFoodNutrientRepository(db)

	// Initialize service with all repositories
	nutrientService := services.NewNutrientService(
		nutrientRepo,
		mealLogRepository,
		mealLogItemsRepository,
		foodRepository,
		foodNutrientsRepository,
	)

	// Initialize controller
	nutrientController := controllers.NewNutrientController(nutrientService)

	// Create auth middleware
	authMiddleware := auth.NewAuthMiddleware()

	// Setup basic nutrient CRUD routes
	nutrientRoutes := router.Group("/nutrients")
	{
		nutrientRoutes.POST("/", nutrientController.CreateNutrient)
		nutrientRoutes.GET("/", nutrientController.GetAllNutrients)
		nutrientRoutes.GET("/:id", nutrientController.GetNutrient)
		nutrientRoutes.GET("/category/:category", nutrientController.GetNutrientsByCategory)
		nutrientRoutes.PUT("/:id", nutrientController.UpdateNutrient)
		nutrientRoutes.DELETE("/:id", nutrientController.DeleteNutrient)
	}

	// Setup nutrition calculation routes (protected)
	nutritionRoutes := router.Group("/nutrition", authMiddleware.RequireAuth())
	{
		nutritionRoutes.GET("/today", nutrientController.GetUserCurrentNutrition)
		nutritionRoutes.GET("/date/:date", nutrientController.GetUserNutritionByDate)
		nutritionRoutes.GET("/range", nutrientController.GetUserNutritionByDateRange)
		nutritionRoutes.GET("/meal/:mealLogId", nutrientController.GetMealNutrition)
	}
}
