package repository

import (
	"github.com/momokapoolz/caloriesapp/food/models"
	"gorm.io/gorm"
)

// FoodRepository handles all database operations for the Food model
type FoodRepository struct {
	db *gorm.DB
}

// NewFoodRepository creates a new food repository instance
func NewFoodRepository(db *gorm.DB) *FoodRepository {
	return &FoodRepository{db: db}
}

// Create adds a new food record to the database
func (r *FoodRepository) Create(food *models.Food) error {
	return r.db.Create(food).Error
}

// GetByID retrieves a food by its ID
func (r *FoodRepository) GetByID(id uint) (*models.Food, error) {
	var food models.Food
	err := r.db.Where("id = ?", id).First(&food).Error
	if err != nil {
		return nil, err
	}
	return &food, nil
}

// GetAll retrieves all foods
func (r *FoodRepository) GetAll() ([]models.Food, error) {
	var foods []models.Food
	err := r.db.Find(&foods).Error
	return foods, err
}

// Update updates a food record
func (r *FoodRepository) Update(food *models.Food) error {
	return r.db.Save(food).Error
}

// Delete removes a food record
func (r *FoodRepository) Delete(id uint) error {
	return r.db.Delete(&models.Food{}, id).Error
} 