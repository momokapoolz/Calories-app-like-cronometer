package services

import (
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