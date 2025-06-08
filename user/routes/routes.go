package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/user/controllers"
	"github.com/momokapoolz/caloriesapp/user/middleware"
	"github.com/momokapoolz/caloriesapp/user/repository"
	"github.com/momokapoolz/caloriesapp/user/services"
)

// SetupRoutes configures all user API routes on the provided router group
func SetupRoutes(rg *gin.RouterGroup) {
	// Apply middleware specific to user routes if needed
	rg.Use(middleware.LoggingMiddleware())
	rg.Use(auth.CORSMiddleware())

	// Set up repositories and services
	userRepo := repository.NewUserRepository()
	passwordService := services.NewPasswordService(userRepo)

	// Set up controllers
	userController := controllers.NewUserController()
	authController := controllers.NewUserAuthController()
	passwordController := controllers.NewPasswordController(passwordService)

	// Register basic routes
	userController.RegisterRoutes(rg)
	authController.RegisterRoutes(rg)

	// Protected routes
	authMiddleware := auth.NewAuthMiddleware()

	// User routes that require authentication
	protected := rg.Group("/user")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.POST("/password/update", passwordController.UpdatePassword)
	}

	// Admin routes
	admin := rg.Group("/admin")
	admin.Use(authMiddleware.RequireAuth())
	admin.Use(authMiddleware.RequireRole("admin"))
	{
		admin.POST("/user/password/update", passwordController.AdminUpdatePassword)
	}

	// Set up JWT auth routes and middleware
	//auth.SetupAuth(rg, userRepo)
}
