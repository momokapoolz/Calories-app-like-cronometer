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
	}
} 