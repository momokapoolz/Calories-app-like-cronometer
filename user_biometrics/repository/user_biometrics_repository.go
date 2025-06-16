package repository

import (
	"time"

	"github.com/momokapoolz/caloriesapp/user_biometrics/models"
	"gorm.io/gorm"
)

// UserBiometricRepository handles all database operations for the UserBiometric model
type UserBiometricRepository struct {
	db *gorm.DB
}

// NewUserBiometricRepository creates a new user biometric repository instance
func NewUserBiometricRepository(db *gorm.DB) *UserBiometricRepository {
	return &UserBiometricRepository{db: db}
}

// Create adds a new user biometric record to the database
func (r *UserBiometricRepository) Create(biometric *models.UserBiometric) error {
	return r.db.Create(biometric).Error
}

// GetByID retrieves a user biometric by its ID
func (r *UserBiometricRepository) GetByID(id uint) (*models.UserBiometric, error) {
	var biometric models.UserBiometric
	err := r.db.Where("id = ?", id).First(&biometric).Error
	if err != nil {
		return nil, err
	}
	return &biometric, nil
}

// GetByUserID retrieves all biometrics for a specific user
func (r *UserBiometricRepository) GetByUserID(userID uint) ([]models.UserBiometric, error) {
	var biometrics []models.UserBiometric
	err := r.db.Where("user_id = ?", userID).Find(&biometrics).Error
	return biometrics, err
}

// GetByUserIDAndType retrieves biometrics of a specific type for a specific user
func (r *UserBiometricRepository) GetByUserIDAndType(userID uint, biometricType string) ([]models.UserBiometric, error) {
	var biometrics []models.UserBiometric
	err := r.db.Where("user_id = ? AND type = ?", userID, biometricType).Order("created_at DESC").Find(&biometrics).Error
	return biometrics, err
}

// GetByUserIDAndTypeAndDateRange retrieves biometrics of a specific type for a user within a date range
func (r *UserBiometricRepository) GetByUserIDAndTypeAndDateRange(userID uint, biometricType string, startDate, endDate time.Time) ([]models.UserBiometric, error) {
	var biometrics []models.UserBiometric
	err := r.db.Where("user_id = ? AND type = ? AND created_at >= ? AND created_at <= ?", userID, biometricType, startDate, endDate).Order("created_at ASC").Find(&biometrics).Error
	return biometrics, err
}

// GetLatestByUserIDAndType retrieves the most recent biometric of a specific type for a user
func (r *UserBiometricRepository) GetLatestByUserIDAndType(userID uint, biometricType string) (*models.UserBiometric, error) {
	var biometric models.UserBiometric
	err := r.db.Where("user_id = ? AND type = ?", userID, biometricType).Order("created_at DESC").First(&biometric).Error
	if err != nil {
		return nil, err
	}
	return &biometric, nil
}

// Update updates a user biometric record
func (r *UserBiometricRepository) Update(biometric *models.UserBiometric) error {
	return r.db.Save(biometric).Error
}

// Delete removes a user biometric record
func (r *UserBiometricRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserBiometric{}, id).Error
}

// GetLatestBiometricsByUserID retrieves the latest biometric for each type for a user
func (r *UserBiometricRepository) GetLatestBiometricsByUserID(userID uint) (map[string]models.UserBiometric, error) {
	//var biometrics []models.UserBiometric
	result := make(map[string]models.UserBiometric)

	// Get all distinct types for the user
	var types []string
	err := r.db.Model(&models.UserBiometric{}).Where("user_id = ?", userID).Distinct("type").Pluck("type", &types).Error
	if err != nil {
		return nil, err
	}

	// Get the latest biometric for each type
	for _, biometricType := range types {
		latest, err := r.GetLatestByUserIDAndType(userID, biometricType)
		if err == nil {
			result[biometricType] = *latest
		}
	}

	return result, nil
}

// GetBiometricTypesForUser retrieves all biometric types that a user has data for
func (r *UserBiometricRepository) GetBiometricTypesForUser(userID uint) ([]string, error) {
	var types []string
	err := r.db.Model(&models.UserBiometric{}).Where("user_id = ?", userID).Distinct("type").Pluck("type", &types).Error
	return types, err
}

// GetProgressDataByType retrieves biometric data for progress tracking within a date range
func (r *UserBiometricRepository) GetProgressDataByType(userID uint, biometricType string, startDate, endDate time.Time, limit int) ([]models.UserBiometric, error) {
	var biometrics []models.UserBiometric
	query := r.db.Where("user_id = ? AND type = ? AND created_at >= ? AND created_at <= ?", userID, biometricType, startDate, endDate).Order("created_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&biometrics).Error
	return biometrics, err
}

// GetBiometricStatistics retrieves statistics for a biometric type
func (r *UserBiometricRepository) GetBiometricStatistics(userID uint, biometricType string, startDate, endDate time.Time) (map[string]interface{}, error) {
	var result struct {
		Count int     `json:"count"`
		Min   float64 `json:"min"`
		Max   float64 `json:"max"`
		Avg   float64 `json:"avg"`
		First float64 `json:"first"`
		Last  float64 `json:"last"`
	}

	query := `
		SELECT 
			COUNT(*) as count,
			MIN(value) as min,
			MAX(value) as max,
			AVG(value) as avg,
			FIRST_VALUE(value) OVER (ORDER BY created_at ASC) as first,
			LAST_VALUE(value) OVER (ORDER BY created_at ASC ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING) as last
		FROM user_biometrics 
		WHERE user_id = ? AND type = ? AND created_at >= ? AND created_at <= ?
		LIMIT 1
	`

	err := r.db.Raw(query, userID, biometricType, startDate, endDate).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"count":  result.Count,
		"min":    result.Min,
		"max":    result.Max,
		"avg":    result.Avg,
		"first":  result.First,
		"last":   result.Last,
		"change": result.Last - result.First,
	}

	if result.First != 0 {
		stats["percent_change"] = ((result.Last - result.First) / result.First) * 100
	}

	return stats, nil
}

// GetBiometricsForAdvancedMetrics retrieves specific biometrics needed for advanced calculations
func (r *UserBiometricRepository) GetBiometricsForAdvancedMetrics(userID uint) (map[string]float64, error) {
	biometricTypes := []string{"weight", "height", "body_fat_percentage", "muscle_mass", "waist_circumference", "hip_circumference", "body_water_percentage"}
	result := make(map[string]float64)

	for _, biometricType := range biometricTypes {
		latest, err := r.GetLatestByUserIDAndType(userID, biometricType)
		if err == nil {
			result[biometricType] = latest.Value
		}
	}

	return result, nil
}
