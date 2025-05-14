package repository

import (
	"github.com/momokapoolz/caloriesapp/food_nutrients/models"
	"gorm.io/gorm"
)

// FoodNutrientRepository handles all database operations for the FoodNutrient model
type FoodNutrientRepository struct {
	db *gorm.DB
}

// NewFoodNutrientRepository creates a new food nutrient repository instance
func NewFoodNutrientRepository(db *gorm.DB) *FoodNutrientRepository {
	return &FoodNutrientRepository{db: db}
}

// Create adds a new food nutrient record to the database
func (r *FoodNutrientRepository) Create(foodNutrient *models.FoodNutrient) error {
	return r.db.Create(foodNutrient).Error
}

// GetByID retrieves a food nutrient by its ID
func (r *FoodNutrientRepository) GetByID(id uint) (*models.FoodNutrient, error) {
	var foodNutrient models.FoodNutrient
	err := r.db.Where("id = ?", id).First(&foodNutrient).Error
	if err != nil {
		return nil, err
	}
	return &foodNutrient, nil
}

// GetAll retrieves all food nutrients
func (r *FoodNutrientRepository) GetAll() ([]models.FoodNutrient, error) {
	var foodNutrients []models.FoodNutrient
	err := r.db.Find(&foodNutrients).Error
	return foodNutrients, err
}

// GetByFoodID retrieves food nutrients by food ID
func (r *FoodNutrientRepository) GetByFoodID(foodID uint) ([]models.FoodNutrient, error) {
	var foodNutrients []models.FoodNutrient
	err := r.db.Where("food_id = ?", foodID).Find(&foodNutrients).Error
	return foodNutrients, err
}

// GetByNutrientID retrieves food nutrients by nutrient ID
func (r *FoodNutrientRepository) GetByNutrientID(nutrientID uint) ([]models.FoodNutrient, error) {
	var foodNutrients []models.FoodNutrient
	err := r.db.Where("nutrient_id = ?", nutrientID).Find(&foodNutrients).Error
	return foodNutrients, err
}

// Update updates a food nutrient record
func (r *FoodNutrientRepository) Update(foodNutrient *models.FoodNutrient) error {
	return r.db.Save(foodNutrient).Error
}

// Delete removes a food nutrient record
func (r *FoodNutrientRepository) Delete(id uint) error {
	return r.db.Delete(&models.FoodNutrient{}, id).Error
} 