package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	dashboard_routes "github.com/momokapoolz/caloriesapp/dashboard/routes"
	"github.com/momokapoolz/caloriesapp/food/routes"
	food_nutrients_routes "github.com/momokapoolz/caloriesapp/food_nutrients/routes"
	meal_log_routes "github.com/momokapoolz/caloriesapp/meal_log/routes"
	meal_log_items_routes "github.com/momokapoolz/caloriesapp/meal_log_items/routes"
	nutrient_routes "github.com/momokapoolz/caloriesapp/nutrient/routes"
	user_biometrics_routes "github.com/momokapoolz/caloriesapp/user_biometrics/routes"

	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(auth.CORSMiddleware())

	v1 := router.Group("/api/v1")

	routes.SetupFoodRoutes(v1, db)
	nutrient_routes.SetupNutrientRoutes(v1, db)
	food_nutrients_routes.SetupFoodNutrientRoutes(v1, db)
	meal_log_routes.SetupMealLogRoutes(v1, db)
	meal_log_items_routes.SetupMealLogItemRoutes(v1, db)
	user_biometrics_routes.SetupUserBiometricRoutes(v1, db)
	dashboard_routes.SetupDashboardRoutes(v1, db)

	return router
}
