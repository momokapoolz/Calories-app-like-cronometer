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
	err := r.db.Where("user_id = ? AND type = ?", userID, biometricType).Find(&biometrics).Error
	return biometrics, err
}

// GetByUserIDAndTypeAndDateRange retrieves biometrics of a specific type for a user within a date range
func (r *UserBiometricRepository) GetByUserIDAndTypeAndDateRange(userID uint, biometricType string, startDate, endDate time.Time) ([]models.UserBiometric, error) {
	var biometrics []models.UserBiometric
	err := r.db.Where("user_id = ? AND type = ? AND created_at >= ? AND created_at <= ?", userID, biometricType, startDate, endDate).Find(&biometrics).Error
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