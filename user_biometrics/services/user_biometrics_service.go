package services

import (
	"math"
	"time"

	"github.com/momokapoolz/caloriesapp/user_biometrics/models"
	"github.com/momokapoolz/caloriesapp/user_biometrics/repository"
)

// UserBiometricService handles business logic for user biometric operations
type UserBiometricService struct {
	repo *repository.UserBiometricRepository
}

// NewUserBiometricService creates a new user biometric service instance
func NewUserBiometricService(repo *repository.UserBiometricRepository) *UserBiometricService {
	return &UserBiometricService{repo: repo}
}

// CreateUserBiometric creates a new user biometric record
func (s *UserBiometricService) CreateUserBiometric(biometric *models.UserBiometric) error {
	return s.repo.Create(biometric)
}

// GetUserBiometricByID retrieves a user biometric record by ID
func (s *UserBiometricService) GetUserBiometricByID(id uint) (*models.UserBiometric, error) {
	return s.repo.GetByID(id)
}

// GetUserBiometricsByUserID retrieves all biometrics for a specific user
func (s *UserBiometricService) GetUserBiometricsByUserID(userID uint) ([]models.UserBiometric, error) {
	return s.repo.GetByUserID(userID)
}

// GetUserBiometricsByUserIDAndType retrieves biometrics of a specific type for a specific user
func (s *UserBiometricService) GetUserBiometricsByUserIDAndType(userID uint, biometricType string) ([]models.UserBiometric, error) {
	return s.repo.GetByUserIDAndType(userID, biometricType)
}

// GetUserBiometricsByUserIDAndTypeAndDateRange retrieves biometrics of a specific type for a user within a date range
func (s *UserBiometricService) GetUserBiometricsByUserIDAndTypeAndDateRange(userID uint, biometricType string, startDate, endDate time.Time) ([]models.UserBiometric, error) {
	return s.repo.GetByUserIDAndTypeAndDateRange(userID, biometricType, startDate, endDate)
}

// GetLatestUserBiometricByUserIDAndType retrieves the most recent biometric of a specific type for a user
func (s *UserBiometricService) GetLatestUserBiometricByUserIDAndType(userID uint, biometricType string) (*models.UserBiometric, error) {
	return s.repo.GetLatestByUserIDAndType(userID, biometricType)
}

// UpdateUserBiometric updates a user biometric record
func (s *UserBiometricService) UpdateUserBiometric(biometric *models.UserBiometric) error {
	return s.repo.Update(biometric)
}

// DeleteUserBiometric removes a user biometric record
func (s *UserBiometricService) DeleteUserBiometric(id uint) error {
	return s.repo.Delete(id)
}

// GetBiometricProgress calculates progress data for a specific biometric type
func (s *UserBiometricService) GetBiometricProgress(userID uint, biometricType string, startDate, endDate time.Time) (*models.BiometricProgress, error) {
	biometrics, err := s.repo.GetByUserIDAndTypeAndDateRange(userID, biometricType, startDate, endDate)
	if err != nil || len(biometrics) == 0 {
		return nil, err
	}

	progress := &models.BiometricProgress{
		Type:      biometricType,
		Unit:      biometrics[0].Unit,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Calculate progress data points
	var dataPoints []models.ProgressData
	for i, biometric := range biometrics {
		var change float64
		var trend string

		if i > 0 {
			change = biometric.Value - biometrics[i-1].Value
			if change > 0 {
				trend = "up"
			} else if change < 0 {
				trend = "down"
			} else {
				trend = "stable"
			}
		}

		dataPoints = append(dataPoints, models.ProgressData{
			Date:   biometric.CreatedAt,
			Value:  biometric.Value,
			Change: change,
			Trend:  trend,
		})
	}

	progress.DataPoints = dataPoints
	progress.CurrentValue = biometrics[len(biometrics)-1].Value
	progress.PreviousValue = biometrics[0].Value
	progress.OverallChange = progress.CurrentValue - progress.PreviousValue

	if progress.PreviousValue != 0 {
		progress.PercentChange = (progress.OverallChange / progress.PreviousValue) * 100
	}

	if progress.OverallChange > 0 {
		progress.Trend = "up"
	} else if progress.OverallChange < 0 {
		progress.Trend = "down"
	} else {
		progress.Trend = "stable"
	}

	return progress, nil
}

// GetChartData returns data formatted for chart visualization
func (s *UserBiometricService) GetChartData(userID uint, biometricType string, startDate, endDate time.Time, maxPoints int) (*models.ChartData, error) {
	biometrics, err := s.repo.GetProgressDataByType(userID, biometricType, startDate, endDate, maxPoints)
	if err != nil || len(biometrics) == 0 {
		return nil, err
	}

	chartData := &models.ChartData{
		Type:      biometricType,
		Unit:      biometrics[0].Unit,
		StartDate: startDate,
		EndDate:   endDate,
	}

	var labels []string
	var values []float64

	for _, biometric := range biometrics {
		labels = append(labels, biometric.CreatedAt.Format("2006-01-02"))
		values = append(values, biometric.Value)
	}

	chartData.Labels = labels
	chartData.Values = values

	return chartData, nil
}

// GetAdvancedMetrics calculates advanced health metrics
func (s *UserBiometricService) GetAdvancedMetrics(userID uint) (*models.AdvancedMetrics, error) {
	biometrics, err := s.repo.GetBiometricsForAdvancedMetrics(userID)
	if err != nil {
		return nil, err
	}

	metrics := &models.AdvancedMetrics{}

	// Calculate BMI
	if weight, hasWeight := biometrics["weight"]; hasWeight {
		if height, hasHeight := biometrics["height"]; hasHeight {
			heightInMeters := height / 100 // assuming height is in cm
			metrics.BMI = weight / (heightInMeters * heightInMeters)
			metrics.BMICategory = s.getBMICategory(metrics.BMI)
		}
	}

	// Set body fat percentage if available
	if bodyFat, hasBodyFat := biometrics["body_fat_percentage"]; hasBodyFat {
		metrics.BodyFatPercentage = bodyFat
		metrics.BodyFatCategory = s.getBodyFatCategory(bodyFat)
	}

	// Set muscle mass if available
	if muscleMass, hasMuscleMass := biometrics["muscle_mass"]; hasMuscleMass {
		metrics.MuscleMass = muscleMass
	}

	// Calculate waist-to-hip ratio
	if waist, hasWaist := biometrics["waist_circumference"]; hasWaist {
		if hip, hasHip := biometrics["hip_circumference"]; hasHip {
			metrics.WaistToHipRatio = waist / hip
		}
	}

	// Set body water percentage if available
	if bodyWater, hasBodyWater := biometrics["body_water_percentage"]; hasBodyWater {
		metrics.BodyWaterPercentage = bodyWater
	}

	// Calculate health risk
	metrics.HealthRisk = s.calculateHealthRisk(metrics)

	return metrics, nil
}

// GetBiometricSummary provides a comprehensive summary of user's biometrics
func (s *UserBiometricService) GetBiometricSummary(userID uint) (*models.BiometricSummary, error) {
	latestBiometrics, err := s.repo.GetLatestBiometricsByUserID(userID)
	if err != nil {
		return nil, err
	}

	summary := &models.BiometricSummary{
		UserID:           userID,
		LatestBiometrics: latestBiometrics,
		ProgressData:     make(map[string]models.BiometricProgress),
		LastUpdated:      time.Now(),
	}

	// Calculate progress for the last 30 days for each biometric type
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	for biometricType := range latestBiometrics {
		progress, err := s.GetBiometricProgress(userID, biometricType, startDate, endDate)
		if err == nil && progress != nil {
			summary.ProgressData[biometricType] = *progress
		}
	}

	return summary, nil
}

// GetAvailableBiometricTypes returns all biometric types available for a user
func (s *UserBiometricService) GetAvailableBiometricTypes(userID uint) ([]string, error) {
	return s.repo.GetBiometricTypesForUser(userID)
}

// Helper functions

func (s *UserBiometricService) getBMICategory(bmi float64) string {
	if bmi < 18.5 {
		return "Underweight"
	} else if bmi < 25 {
		return "Normal weight"
	} else if bmi < 30 {
		return "Overweight"
	} else {
		return "Obese"
	}
}

func (s *UserBiometricService) getBodyFatCategory(bodyFat float64) string {
	// These are general categories and may vary by age and gender
	if bodyFat < 10 {
		return "Essential fat"
	} else if bodyFat < 14 {
		return "Athletes"
	} else if bodyFat < 21 {
		return "Fitness"
	} else if bodyFat < 25 {
		return "Average"
	} else {
		return "Obese"
	}
}

func (s *UserBiometricService) calculateHealthRisk(metrics *models.AdvancedMetrics) string {
	risk := "Low"

	if metrics.BMI >= 30 {
		risk = "High"
	} else if metrics.BMI >= 25 {
		risk = "Moderate"
	}

	if metrics.WaistToHipRatio > 0.95 { // For men, adjust for women
		risk = "High"
	} else if metrics.WaistToHipRatio > 0.90 {
		if risk == "Low" {
			risk = "Moderate"
		}
	}

	return risk
}

// CalculateTrend determines the trend direction based on data points
func (s *UserBiometricService) CalculateTrend(values []float64) string {
	if len(values) < 2 {
		return "stable"
	}

	// Simple linear regression to determine trend
	n := float64(len(values))
	sumX, sumY, sumXY, sumXX := 0.0, 0.0, 0.0, 0.0

	for i, y := range values {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)

	if math.Abs(slope) < 0.001 {
		return "stable"
	} else if slope > 0 {
		return "up"
	} else {
		return "down"
	}
}
