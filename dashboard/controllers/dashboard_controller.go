package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/dashboard/services"
	"github.com/momokapoolz/caloriesapp/helpers"
)

type DashboardController struct {
	service *services.DashboardService
}

// NewDashboardController creates a new dashboard controller instance
func NewDashboardController(service *services.DashboardService) *DashboardController {
	return &DashboardController{service: service}
}

// GetUserDashboard godoc
// @Summary      Get user dashboard
// @Description  Get calorie summary and nutrition data for a specific date
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Param        date  query     string  false  "Date in YYYY-MM-DD format (default: today)"
// @Success      200   {object}  dto.DashboardResponseDTO  "Dashboard data retrieved successfully"
// @Failure      400   {object}  map[string]string         "Invalid date format"
// @Failure      401   {object}  map[string]string         "Unauthorized"
// @Failure      500   {object}  map[string]string         "Internal server error"
// @Security     BearerAuth
// @Router       /dashboard/ [get]
func (c *DashboardController) GetUserDashboard(ctx *gin.Context) {
	//Get current user from JWT token
	userClaims, ok := auth.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//Get date from query parameter (if not have, default is today)
	dateStr := ctx.DefaultQuery("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	//Get dashboard data from service
	dashboard, err := c.service.GetUserDashboard(userClaims.UserID, date)
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard data: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dashboard)
}
