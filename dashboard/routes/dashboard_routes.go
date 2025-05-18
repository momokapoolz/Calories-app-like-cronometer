package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/dashboard/controllers"
	"github.com/momokapoolz/caloriesapp/dashboard/services"
	foodRepository "github.com/momokapoolz/caloriesapp/food/repository"
	foodNutrientsRepository "github.com/momokapoolz/caloriesapp/food_nutrients/repository"
	mealLogRepository "github.com/momokapoolz/caloriesapp/meal_log/repository"
	mealLogItemsRepository "github.com/momokapoolz/caloriesapp/meal_log_items/repository"
	nutrientRepository "github.com/momokapoolz/caloriesapp/nutrient/repository"
	"gorm.io/gorm"
)

// SetupDashboardRoutes initializes dashboard routes
func SetupDashboardRoutes(router *gin.RouterGroup, db *gorm.DB) {
	// Initialize repositories
	mealLogRepo := mealLogRepository.NewMealLogRepository(db)
	mealLogItemsRepo := mealLogItemsRepository.NewMealLogItemRepository(db)
	foodRepo := foodRepository.NewFoodRepository(db)
	nutrientRepo := nutrientRepository.NewNutrientRepository(db)
	foodNutrientsRepo := foodNutrientsRepository.NewFoodNutrientRepository(db)

	// Initialize service
	dashboardService := services.NewDashboardService(
		mealLogRepo,
		mealLogItemsRepo,
		foodRepo,
		nutrientRepo,
		foodNutrientsRepo,
	)

	// Initialize controller
	dashboardController := controllers.NewDashboardController(dashboardService)

	// Create auth middleware
	authMiddleware := auth.NewAuthMiddleware()

	// Set up protected routes
	dashboardRoutes := router.Group("/dashboard", authMiddleware.RequireAuth())
	{
		dashboardRoutes.GET("/", dashboardController.GetUserDashboard)
	}
}
