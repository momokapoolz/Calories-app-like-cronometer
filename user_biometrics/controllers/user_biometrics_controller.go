package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/helpers"
	"github.com/momokapoolz/caloriesapp/user_biometrics/models"
	"github.com/momokapoolz/caloriesapp/user_biometrics/services"
)

type UserBiometricController struct {
	service *services.UserBiometricService
}

// NewUserBiometricController creates a new user biometric controller instance
func NewUserBiometricController(service *services.UserBiometricService) *UserBiometricController {
	return &UserBiometricController{service: service}
}

// CreateUserBiometric godoc
// @Summary      Create user biometric
// @Description  Create a new user biometric record
// @Tags         user_biometric
// @Accept       json
// @Produce      json
// @Param        biometric  body      models.UserBiometric  true  "User biometric data"
// @Success      201  {object}  models.UserBiometric  "User biometric created successfully"
// @Failure      400  {object}  map[string]string     "Invalid request body"
// @Failure      500  {object}  map[string]string     "Internal server error"
// @Security     BearerAuth
// @Router       /user-biometrics/ [post]
func (c *UserBiometricController) CreateUserBiometric(ctx *gin.Context) {
	var biometric models.UserBiometric
	if err := ctx.ShouldBindJSON(&biometric); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set created_at to current time if not provided
	if biometric.CreatedAt.IsZero() {
		biometric.CreatedAt = time.Now()
	}

	if err := c.service.CreateUserBiometric(&biometric); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user biometric"})
		return
	}

	ctx.JSON(http.StatusCreated, biometric)
}

// GetUserBiometric godoc
// @Summary      Get user biometric by ID
// @Description  Retrieve a user biometric record by its ID
// @Tags         user_biometric
// @Produce      json
// @Param        id  path      int  true  "Biometric record ID"
// @Success      200  {object}  models.UserBiometric  "User biometric retrieved successfully"
// @Failure      400  {object}  map[string]string     "Invalid ID format"
// @Failure      404  {object}  map[string]string     "User biometric not found"
// @Security     BearerAuth
// @Router       /user-biometrics/{id} [get]
func (c *UserBiometricController) GetUserBiometric(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	biometric, err := c.service.GetUserBiometricByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User biometric not found"})
		return
	}

	ctx.JSON(http.StatusOK, biometric)
}

// GetUserBiometricsByUserID godoc
// @Summary      Get all biometrics for a user
// @Description  Retrieve all biometric records for a specific user
// @Tags         user_biometric
// @Produce      json
// @Param        userId  path      int  true  "User ID"
// @Success      200  {array}   models.UserBiometric  "List of user biometrics"
// @Failure      400  {object}  map[string]string     "Invalid user ID format"
// @Failure      500  {object}  map[string]string     "Internal server error"
// @Security     BearerAuth
// @Router       /user-biometrics/user/{userId} [get]
func (c *UserBiometricController) GetUserBiometricsByUserID(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	biometrics, err := c.service.GetUserBiometricsByUserID(uint(userID))
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user biometrics"})
		return
	}

	ctx.JSON(http.StatusOK, biometrics)
}

// GetUserBiometricsByUserIDAndType godoc
// @Summary      Get biometrics by user and type
// @Description  Retrieve biometric records of a specific type for a specific user
// @Tags         user_biometric
// @Produce      json
// @Param        userId  path      int     true  "User ID"
// @Param        type    path      string  true  "Biometric type (e.g. weight, height, bmi)"
// @Success      200  {array}   models.UserBiometric  "List of user biometrics by type"
// @Failure      400  {object}  map[string]string     "Invalid user ID format"
// @Failure      500  {object}  map[string]string     "Internal server error"
// @Security     BearerAuth
// @Router       /user-biometrics/user/{userId}/type/{type} [get]
func (c *UserBiometricController) GetUserBiometricsByUserIDAndType(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	biometricType := ctx.Param("type")
	biometrics, err := c.service.GetUserBiometricsByUserIDAndType(uint(userID), biometricType)
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user biometrics by type"})
		return
	}

	ctx.JSON(http.StatusOK, biometrics)
}

// GetUserBiometricsByUserIDAndTypeAndDateRange godoc
// @Summary      Get biometrics by user, type and date range
// @Description  Retrieve biometric records of a specific type for a user within a date range
// @Tags         user_biometric
// @Produce      json
// @Param        userId     path   int     true  "User ID"
// @Param        type       path   string  true  "Biometric type"
// @Param        startDate  query  string  true  "Start date in YYYY-MM-DD format"
// @Param        endDate    query  string  true  "End date in YYYY-MM-DD format"
// @Success      200  {array}   models.UserBiometric  "List of user biometrics in date range"
// @Failure      400  {object}  map[string]string     "Invalid parameters or date format"
// @Failure      500  {object}  map[string]string     "Internal server error"
// @Security     BearerAuth
// @Router       /user-biometrics/user/{userId}/type/{type}/date-range [get]
func (c *UserBiometricController) GetUserBiometricsByUserIDAndTypeAndDateRange(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	biometricType := ctx.Param("type")
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Start date and end date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Set end date to the end of the day
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	biometrics, err := c.service.GetUserBiometricsByUserIDAndTypeAndDateRange(uint(userID), biometricType, startDate, endDate)
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user biometrics for the specified date range"})
		return
	}

	ctx.JSON(http.StatusOK, biometrics)
}

// GetLatestUserBiometricByUserIDAndType godoc
// @Summary      Get latest biometric by user and type
// @Description  Retrieve the most recent biometric record of a specific type for a user
// @Tags         user_biometric
// @Produce      json
// @Param        userId  path  int     true  "User ID"
// @Param        type    path  string  true  "Biometric type"
// @Success      200  {object}  models.UserBiometric  "Latest biometric record"
// @Failure      400  {object}  map[string]string     "Invalid user ID format"
// @Failure      404  {object}  map[string]string     "Latest user biometric not found"
// @Security     BearerAuth
// @Router       /user-biometrics/user/{userId}/type/{type}/latest [get]
func (c *UserBiometricController) GetLatestUserBiometricByUserIDAndType(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	biometricType := ctx.Param("type")
	biometric, err := c.service.GetLatestUserBiometricByUserIDAndType(uint(userID), biometricType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Latest user biometric not found"})
		return
	}

	ctx.JSON(http.StatusOK, biometric)
}

// UpdateUserBiometric godoc
// @Summary      Update user biometric
// @Description  Update an existing user biometric record by ID
// @Tags         user_biometric
// @Accept       json
// @Produce      json
// @Param        id         path  int                   true  "Biometric record ID"
// @Param        biometric  body  models.UserBiometric  true  "Updated biometric data"
// @Success      200  {object}  models.UserBiometric  "User biometric updated successfully"
// @Failure      400  {object}  map[string]string     "Invalid ID or request body"
// @Failure      500  {object}  map[string]string     "Internal server error"
// @Security     BearerAuth
// @Router       /user-biometrics/{id} [put]
func (c *UserBiometricController) UpdateUserBiometric(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var biometric models.UserBiometric
	if err := ctx.ShouldBindJSON(&biometric); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	biometric.ID = uint(id)
	if err := c.service.UpdateUserBiometric(&biometric); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user biometric"})
		return
	}

	ctx.JSON(http.StatusOK, biometric)
}

// DeleteUserBiometric godoc
// @Summary      Delete user biometric
// @Description  Remove a user biometric record by ID
// @Tags         user_biometric
// @Produce      json
// @Param        id  path  int  true  "Biometric record ID"
// @Success      200  {object}  map[string]string  "User biometric deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid ID format"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /user-biometrics/{id} [delete]
func (c *UserBiometricController) DeleteUserBiometric(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.DeleteUserBiometric(uint(id)); err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user biometric"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User biometric deleted successfully"})
}

// GetBiometricProgress godoc
// @Summary      Get biometric progress
// @Description  Retrieve progress data for a specific biometric type over a date range (defaults to last 30 days)
// @Tags         user_biometric
// @Produce      json
// @Param        userId     path   int     true   "User ID"
// @Param        type       path   string  true   "Biometric type"
// @Param        startDate  query  string  false  "Start date in YYYY-MM-DD format (default: 30 days ago)"
// @Param        endDate    query  string  false  "End date in YYYY-MM-DD format (default: today)"
// @Success      200  {object}  map[string]interface{}  "Biometric progress data"
// @Failure      400  {object}  map[string]string       "Invalid user ID format"
// @Failure      500  {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /user-biometrics/user/{userId}/progress/{type} [get]
func (c *UserBiometricController) GetBiometricProgress(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	biometricType := ctx.Param("type")
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")

	// Default to last 30 days if no dates provided
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	if startDateStr != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsedStartDate
		}
	}

	if endDateStr != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = time.Date(parsedEndDate.Year(), parsedEndDate.Month(), parsedEndDate.Day(), 23, 59, 59, 999999999, parsedEndDate.Location())
		}
	}

	progress, err := c.service.GetBiometricProgress(uint(userID), biometricType, startDate, endDate)
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve biometric progress"})
		return
	}

	ctx.JSON(http.StatusOK, progress)
}

// GetChartData godoc
// @Summary      Get biometric chart data
// @Description  Retrieve biometric data formatted for chart visualization over a date range
// @Tags         user_biometric
// @Produce      json
// @Param        userId     path   int     true   "User ID"
// @Param        type       path   string  true   "Biometric type"
// @Param        startDate  query  string  false  "Start date in YYYY-MM-DD format (default: 30 days ago)"
// @Param        endDate    query  string  false  "End date in YYYY-MM-DD format (default: today)"
// @Param        maxPoints  query  int     false  "Maximum number of data points to return (default: 50)"
// @Success      200  {object}  map[string]interface{}  "Chart data"
// @Failure      400  {object}  map[string]string       "Invalid user ID format"
// @Failure      500  {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /user-biometrics/user/{userId}/chart/{type} [get]
func (c *UserBiometricController) GetChartData(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	biometricType := ctx.Param("type")
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")
	maxPointsStr := ctx.DefaultQuery("maxPoints", "50")

	maxPoints, err := strconv.Atoi(maxPointsStr)
	if err != nil {
		maxPoints = 50
	}

	// Default to last 30 days if no dates provided
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	if startDateStr != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsedStartDate
		}
	}

	if endDateStr != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = time.Date(parsedEndDate.Year(), parsedEndDate.Month(), parsedEndDate.Day(), 23, 59, 59, 999999999, parsedEndDate.Location())
		}
	}

	chartData, err := c.service.GetChartData(uint(userID), biometricType, startDate, endDate, maxPoints)
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chart data"})
		return
	}

	ctx.JSON(http.StatusOK, chartData)
}

// GetAdvancedMetrics godoc
// @Summary      Get advanced health metrics for a user
// @Description  Calculate and retrieve advanced health metrics (BMI, body fat, waist-hip ratio, health risk) for a user
// @Tags         user_biometric
// @Produce      json
// @Param        userId  path  int  true  "User ID"
// @Success      200  {object}  map[string]interface{}  "Advanced health metrics"
// @Failure      400  {object}  map[string]string       "Invalid user ID format"
// @Failure      500  {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /user-biometrics/user/{userId}/advanced-metrics [get]
func (c *UserBiometricController) GetAdvancedMetrics(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	metrics, err := c.service.GetAdvancedMetrics(uint(userID))
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate advanced metrics"})
		return
	}

	ctx.JSON(http.StatusOK, metrics)
}

// GetBiometricSummary godoc
// @Summary      Get biometric summary for a user
// @Description  Retrieve a comprehensive summary of all biometric data for a specific user
// @Tags         user_biometric
// @Produce      json
// @Param        userId  path  int  true  "User ID"
// @Success      200  {object}  map[string]interface{}  "Biometric summary"
// @Failure      400  {object}  map[string]string       "Invalid user ID format"
// @Failure      500  {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /user-biometrics/user/{userId}/summary [get]
func (c *UserBiometricController) GetBiometricSummary(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	summary, err := c.service.GetBiometricSummary(uint(userID))
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve biometric summary"})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// GetAvailableBiometricTypes godoc
// @Summary      Get available biometric types for a user
// @Description  Retrieve all biometric types that have recorded data for a specific user
// @Tags         user_biometric
// @Produce      json
// @Param        userId  path  int  true  "User ID"
// @Success      200  {object}  map[string]interface{}  "Available biometric types"
// @Failure      400  {object}  map[string]string       "Invalid user ID format"
// @Failure      500  {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /user-biometrics/user/{userId}/types [get]
func (c *UserBiometricController) GetAvailableBiometricTypes(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	types, err := c.service.GetAvailableBiometricTypes(uint(userID))
	if err != nil {
		helpers.LogError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve available biometric types"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"types": types})
}

// GetBiometricTypes godoc
// @Summary      Get all supported biometric types
// @Description  Retrieve the full list of biometric types supported by the system
// @Tags         user_biometric
// @Produce      json
// @Success      200  {array}   string  "List of supported biometric types"
// @Security     BearerAuth
// @Router       /user-biometrics/types [get]
func (c *UserBiometricController) GetBiometricTypes(ctx *gin.Context) {
	types := models.GetBiometricTypes()
	ctx.JSON(http.StatusOK, types)
}
