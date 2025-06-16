package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/user_biometrics/controllers"
	"github.com/momokapoolz/caloriesapp/user_biometrics/repository"
	"github.com/momokapoolz/caloriesapp/user_biometrics/services"
	"gorm.io/gorm"
)

// SetupUserBiometricRoutes initializes user biometric routes
func SetupUserBiometricRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userBiometricRepo := repository.NewUserBiometricRepository(db)
	userBiometricService := services.NewUserBiometricService(userBiometricRepo)
	userBiometricController := controllers.NewUserBiometricController(userBiometricService)

	userBiometricRoutes := router.Group("/user-biometrics")
	{
		userBiometricRoutes.POST("/", userBiometricController.CreateUserBiometric)
		userBiometricRoutes.GET("/:id", userBiometricController.GetUserBiometric)
		userBiometricRoutes.GET("/user/:userId", userBiometricController.GetUserBiometricsByUserID)
		userBiometricRoutes.GET("/user/:userId/type/:type", userBiometricController.GetUserBiometricsByUserIDAndType)
		userBiometricRoutes.GET("/user/:userId/type/:type/date-range", userBiometricController.GetUserBiometricsByUserIDAndTypeAndDateRange)
		userBiometricRoutes.GET("/user/:userId/type/:type/latest", userBiometricController.GetLatestUserBiometricByUserIDAndType)
		userBiometricRoutes.PUT("/:id", userBiometricController.UpdateUserBiometric)
		userBiometricRoutes.DELETE("/:id", userBiometricController.DeleteUserBiometric)

		userBiometricRoutes.GET("/user/:userId/progress/:type", userBiometricController.GetBiometricProgress)
		userBiometricRoutes.GET("/user/:userId/chart/:type", userBiometricController.GetChartData)
		userBiometricRoutes.GET("/user/:userId/advanced-metrics", userBiometricController.GetAdvancedMetrics)
		userBiometricRoutes.GET("/user/:userId/summary", userBiometricController.GetBiometricSummary)
		userBiometricRoutes.GET("/user/:userId/types", userBiometricController.GetAvailableBiometricTypes)
		userBiometricRoutes.GET("/types", userBiometricController.GetBiometricTypes)
	}
}
