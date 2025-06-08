package repository

import (
	"github.com/momokapoolz/caloriesapp/nutrient/models"
	"gorm.io/gorm"
)

type NutrientRepository struct {
	db *gorm.DB
}

func NewNutrientRepository(db *gorm.DB) *NutrientRepository {
	return &NutrientRepository{db: db}
}

func (r *NutrientRepository) Create(nutrient *models.Nutrient) error {
	return r.db.Create(nutrient).Error
}

func (r *NutrientRepository) GetByID(id uint) (*models.Nutrient, error) {
	var nutrient models.Nutrient
	err := r.db.Where("id = ?", id).First(&nutrient).Error
	if err != nil {
		return nil, err
	}
	return &nutrient, nil
}

func (r *NutrientRepository) GetAll() ([]models.Nutrient, error) {
	var nutrients []models.Nutrient
	err := r.db.Find(&nutrients).Error
	return nutrients, err
}

func (r *NutrientRepository) GetByCategory(category string) ([]models.Nutrient, error) {
	var nutrients []models.Nutrient
	err := r.db.Where("category = ?", category).Find(&nutrients).Error
	return nutrients, err
}

func (r *NutrientRepository) Update(nutrient *models.Nutrient) error {
	return r.db.Save(nutrient).Error
}

func (r *NutrientRepository) Delete(id uint) error {
	return r.db.Delete(&models.Nutrient{}, id).Error
}
