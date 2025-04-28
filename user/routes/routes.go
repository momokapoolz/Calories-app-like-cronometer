package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/user/controllers"
	"github.com/momokapoolz/caloriesapp/user/middleware"
	"github.com/momokapoolz/caloriesapp/user/repository"
)

// SetupRoutes configures all API routes
func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.LoggingMiddleware())

	// Set up user controller
	userRepo := repository.NewUserRepository()
	userController := controllers.NewUserController()
	userController.RegisterRoutes(router)

	// Set up auth controller for login
	authController := controllers.NewUserAuthController()
	authController.RegisterRoutes(router)

	// Set up JWT auth routes and middleware
	auth.SetupAuth(router, userRepo)

	return router
}
