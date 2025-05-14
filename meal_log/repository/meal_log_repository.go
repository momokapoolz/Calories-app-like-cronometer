package repository

import (
	"time"
	"github.com/momokapoolz/caloriesapp/meal_log/models"
	"gorm.io/gorm"
)

// MealLogRepository handles all database operations for the MealLog model
type MealLogRepository struct {
	db *gorm.DB
}

// NewMealLogRepository creates a new meal log repository instance
func NewMealLogRepository(db *gorm.DB) *MealLogRepository {
	return &MealLogRepository{db: db}
}

// Create adds a new meal log record to the database
func (r *MealLogRepository) Create(mealLog *models.MealLog) error {
	return r.db.Create(mealLog).Error
}

// GetByID retrieves a meal log by its ID
func (r *MealLogRepository) GetByID(id uint) (*models.MealLog, error) {
	var mealLog models.MealLog
	err := r.db.Where("id = ?", id).First(&mealLog).Error
	if err != nil {
		return nil, err
	}
	return &mealLog, nil
}

// GetByUserID retrieves all meal logs for a specific user
func (r *MealLogRepository) GetByUserID(userID uint) ([]models.MealLog, error) {
	var mealLogs []models.MealLog
	err := r.db.Where("user_id = ?", userID).Find(&mealLogs).Error
	return mealLogs, err
}

// GetByUserIDAndDate retrieves meal logs for a specific user on a specific date
func (r *MealLogRepository) GetByUserIDAndDate(userID uint, date time.Time) ([]models.MealLog, error) {
	// Format the date to match the date portion only
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.Add(24 * time.Hour)

	var mealLogs []models.MealLog
	err := r.db.Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startDate, endDate).Find(&mealLogs).Error
	return mealLogs, err
}

// GetByUserIDAndDateRange retrieves meal logs for a specific user within a date range
func (r *MealLogRepository) GetByUserIDAndDateRange(userID uint, startDate, endDate time.Time) ([]models.MealLog, error) {
	var mealLogs []models.MealLog
	err := r.db.Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate).Find(&mealLogs).Error
	return mealLogs, err
}

// Update updates a meal log record
func (r *MealLogRepository) Update(mealLog *models.MealLog) error {
	return r.db.Save(mealLog).Error
}

// Delete removes a meal log record
func (r *MealLogRepository) Delete(id uint) error {
	return r.db.Delete(&models.MealLog{}, id).Error
} 