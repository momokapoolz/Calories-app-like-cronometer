package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/dashboard/services"
)

// DashboardController handles HTTP requests for dashboard operations
type DashboardController struct {
	service *services.DashboardService
}

// NewDashboardController creates a new dashboard controller instance
func NewDashboardController(service *services.DashboardService) *DashboardController {
	return &DashboardController{service: service}
}

// GetUserDashboard retrieves the dashboard data for the currently authenticated user
func (c *DashboardController) GetUserDashboard(ctx *gin.Context) {
	// Get current user from JWT token
	userClaims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get date from query parameters (optional, default to today)
	dateStr := ctx.DefaultQuery("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Get dashboard data from service
	dashboard, err := c.service.GetUserDashboard(userClaims.UserID, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard data: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dashboard)
}
