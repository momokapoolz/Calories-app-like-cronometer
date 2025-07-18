package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/dashboard/services"
	"net/http"
	"time"
)

type DashboardController struct {
	service *services.DashboardService
}

func NewDashboardController(service *services.DashboardService) *DashboardController {
	return &DashboardController{service: service}
}

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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard data: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dashboard)
}
