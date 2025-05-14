package repository

import (
	"github.com/momokapoolz/caloriesapp/nutrient/models"
	"gorm.io/gorm"
)

// NutrientRepository handles all database operations for the Nutrient model
type NutrientRepository struct {
	db *gorm.DB
}

// NewNutrientRepository creates a new nutrient repository instance
func NewNutrientRepository(db *gorm.DB) *NutrientRepository {
	return &NutrientRepository{db: db}
}

// Create adds a new nutrient record to the database
func (r *NutrientRepository) Create(nutrient *models.Nutrient) error {
	return r.db.Create(nutrient).Error
}

// GetByID retrieves a nutrient by its ID
func (r *NutrientRepository) GetByID(id uint) (*models.Nutrient, error) {
	var nutrient models.Nutrient
	err := r.db.Where("id = ?", id).First(&nutrient).Error
	if err != nil {
		return nil, err
	}
	return &nutrient, nil
}

// GetAll retrieves all nutrients
func (r *NutrientRepository) GetAll() ([]models.Nutrient, error) {
	var nutrients []models.Nutrient
	err := r.db.Find(&nutrients).Error
	return nutrients, err
}

// GetByCategory retrieves nutrients by category
func (r *NutrientRepository) GetByCategory(category string) ([]models.Nutrient, error) {
	var nutrients []models.Nutrient
	err := r.db.Where("category = ?", category).Find(&nutrients).Error
	return nutrients, err
}

// Update updates a nutrient record
func (r *NutrientRepository) Update(nutrient *models.Nutrient) error {
	return r.db.Save(nutrient).Error
}

// Delete removes a nutrient record
func (r *NutrientRepository) Delete(id uint) error {
	return r.db.Delete(&models.Nutrient{}, id).Error
} 