package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/user/controllers"
	"github.com/momokapoolz/caloriesapp/user/middleware"
	"github.com/momokapoolz/caloriesapp/user/repository"
	"github.com/momokapoolz/caloriesapp/user/services"
)

// SetupRoutes configures all user-module API routes on the provided router group.
// A single AuthMiddleware instance is created here and shared across all sub-routers
// to avoid redundant instantiation.
func SetupRoutes(rg *gin.RouterGroup) {
	rg.Use(middleware.LoggingMiddleware())
	rg.Use(auth.CORSMiddleware())

	// Single shared instance — prevents creating multiple JWTService objects
	authMiddleware := auth.NewAuthMiddleware()

	// Shared dependencies
	userRepo := repository.NewUserRepository()
	passwordService := services.NewPasswordService(userRepo)

	// Controllers
	userController := controllers.NewUserController()
	passwordController := controllers.NewPasswordController(passwordService)

	// Auth routes: POST /login, /register, /refresh, /logout
	SetupAuthRoutes(rg, authMiddleware)

	// User profile routes: GET /profile, PUT /profile, DELETE /account
	userController.RegisterRoutes(rg, authMiddleware)

	// Password update (for authenticated users)
	userProtected := rg.Group("/user")
	userProtected.Use(authMiddleware.RequireAuth())
	{
		userProtected.POST("/password/update", passwordController.UpdatePassword)
	}

	// Admin-only routes
	admin := rg.Group("/admin")
	admin.Use(authMiddleware.RequireAuth())
	admin.Use(authMiddleware.RequireRole("admin"))
	{
		admin.POST("/user/password/update", passwordController.AdminUpdatePassword)
	}
}
