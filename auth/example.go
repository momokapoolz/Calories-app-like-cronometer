package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/user/repository"
)

// SetupAuth demonstrates how to configure auth in your application
func SetupAuth(router *gin.Engine, userRepo *repository.UserRepository) {
	// Create auth controller
	authController := NewAuthController(userRepo)
	
	// Register auth routes (login, refresh)
	authController.RegisterRoutes(router)
	
	// Create auth middleware
	authMiddleware := NewAuthMiddleware()
	
	// Example of a protected route group - using a different path to avoid conflicts
	protected := router.Group("/api/auth")
	{
		// Apply auth middleware to all routes in this group
		protected.Use(authMiddleware.RequireAuth())
		
		// Routes that require authentication - using a different path
		protected.GET("/profile", GetUserProfile)
		
		// Example of role-based protection
		admin := protected.Group("/admin")
		admin.Use(authMiddleware.RequireRole("admin"))
		
		// Routes that require admin role
		admin.GET("/users", ListAllUsers)
	}
}

// GetUserProfile is an example of accessing the authenticated user in a handler
func GetUserProfile(c *gin.Context) {
	// Get the current authenticated user from context
	userClaims, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(401, gin.H{"error": "Not authenticated"})
		return
	}
	
	// Now you can use userClaims.UserID, userClaims.Email, userClaims.Role
	// to identify the authenticated user making the request
	
	// You would typically load the full user profile from the database here
	c.JSON(200, gin.H{
		"user_id": userClaims.UserID,
		"email":   userClaims.Email,
		"role":    userClaims.Role,
	})
}

// ListAllUsers is an example of a handler requiring admin role
func ListAllUsers(c *gin.Context) {
	// This handler will only be executed if the user has the "admin" role
	// You can implement your admin-only functionality here
	
	c.JSON(200, gin.H{
		"message": "This is an admin-only endpoint",
	})
} 