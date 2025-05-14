package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/user/controllers"
	"github.com/momokapoolz/caloriesapp/user/middleware"
	"github.com/momokapoolz/caloriesapp/user/repository"
)

// SetupRoutes configures all user API routes on the provided router group
func SetupRoutes(rg *gin.RouterGroup) {
	// Apply middleware specific to user routes if needed
	rg.Use(middleware.LoggingMiddleware())

	// Set up user controller
	userRepo := repository.NewUserRepository()
	userController := controllers.NewUserController()
	userController.RegisterRoutes(rg)

	// Set up auth controller for login
	authController := controllers.NewUserAuthController()
	authController.RegisterRoutes(rg)

	// Set up JWT auth routes and middleware
	auth.SetupAuth(rg, userRepo)
}
