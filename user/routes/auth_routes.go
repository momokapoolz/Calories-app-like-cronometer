package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/user/controllers"
)

// SetupAuthRoutes registers all authentication endpoints on the given router group.
// A shared authMiddleware instance is passed in to avoid redundant instantiation.
//
// Public routes:
//
//	POST /login    — authenticate, receive HttpOnly JWT cookies
//	POST /register — create account
//	POST /refresh  — rotate access token using refresh_token cookie
//
// Protected routes:
//
//	POST /logout   — clear both auth cookies
func SetupAuthRoutes(rg *gin.RouterGroup, authMiddleware *auth.AuthMiddleware) {
	authController := controllers.NewUserAuthController()

	rg.POST("/login", authController.Login)
	rg.POST("/register", authController.Register)
	rg.POST("/refresh", authController.Refresh)

	rg.POST("/logout", authMiddleware.RequireAuth(), authController.Logout)
}
