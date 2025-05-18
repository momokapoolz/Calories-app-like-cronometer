package repository

import (
	"errors"

	"github.com/momokapoolz/caloriesapp/meal_log_items/models"
	"gorm.io/gorm"
)

// MealLogItemRepository handles all database operations for the MealLogItem model
type MealLogItemRepository struct {
	db *gorm.DB
}

// NewMealLogItemRepository creates a new meal log item repository instance
func NewMealLogItemRepository(db *gorm.DB) *MealLogItemRepository {
	return &MealLogItemRepository{db: db}
}

// Create adds a new meal log item to the database
func (r *MealLogItemRepository) Create(item *models.MealLogItem) error {
	return r.db.Create(item).Error
}

// GetByID retrieves a meal log item by its ID
func (r *MealLogItemRepository) GetByID(id uint) (*models.MealLogItem, error) {
	var item models.MealLogItem
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetByMealLogID retrieves all items for a specific meal log
func (r *MealLogItemRepository) GetByMealLogID(mealLogID uint) ([]models.MealLogItem, error) {
	var items []models.MealLogItem
	err := r.db.Where("meal_log_id = ?", mealLogID).Find(&items).Error
	return items, err
}

// GetByFoodID retrieves all meal log items for a specific food
func (r *MealLogItemRepository) GetByFoodID(foodID uint) ([]models.MealLogItem, error) {
	var items []models.MealLogItem
	err := r.db.Where("food_id = ?", foodID).Find(&items).Error
	return items, err
}

// Update updates a meal log item
func (r *MealLogItemRepository) Update(item *models.MealLogItem) error {
	return r.db.Save(item).Error
}

// Delete removes a meal log item
func (r *MealLogItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.MealLogItem{}, id).Error
}

// DeleteByMealLogID removes all items for a specific meal log
func (r *MealLogItemRepository) DeleteByMealLogID(mealLogID uint) error {
	return r.db.Where("meal_log_id = ?", mealLogID).Delete(&models.MealLogItem{}).Error
}

// CreateBatch adds multiple meal log items to the database in a single transaction
func (r *MealLogItemRepository) CreateBatch(items []models.MealLogItem) ([]models.MealLogItem, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for i := range items {
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

// VerifyMealLogOwnership checks if the meal log with the given ID belongs to the specified user
func (r *MealLogItemRepository) VerifyMealLogOwnership(mealLogID, userID uint) (bool, error) {
	type MealLog struct {
		ID     uint `gorm:"primaryKey"`
		UserID uint
	}

	var mealLog MealLog
	result := r.db.Table("meal_logs").Select("id", "user_id").Where("id = ?", mealLogID).First(&mealLog)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	// Check if the meal log belongs to the user
	return mealLog.UserID == userID, nil
}
