package services

import (
	"github.com/momokapoolz/caloriesapp/food_nutrients/models"
	"github.com/momokapoolz/caloriesapp/food_nutrients/repository"
)

// FoodNutrientService handles business logic for food nutrient operations
type FoodNutrientService struct {
	repo *repository.FoodNutrientRepository
}

// NewFoodNutrientService creates a new food nutrient service instance
func NewFoodNutrientService(repo *repository.FoodNutrientRepository) *FoodNutrientService {
	return &FoodNutrientService{repo: repo}
}

// CreateFoodNutrient creates a new food nutrient record
func (s *FoodNutrientService) CreateFoodNutrient(foodNutrient *models.FoodNutrient) error {
	return s.repo.Create(foodNutrient)
}

// GetFoodNutrientByID retrieves a food nutrient record by ID
func (s *FoodNutrientService) GetFoodNutrientByID(id uint) (*models.FoodNutrient, error) {
	return s.repo.GetByID(id)
}

// GetAllFoodNutrients retrieves all food nutrient records
func (s *FoodNutrientService) GetAllFoodNutrients() ([]models.FoodNutrient, error) {
	return s.repo.GetAll()
}

// GetFoodNutrientsByFoodID retrieves food nutrients by food ID
func (s *FoodNutrientService) GetFoodNutrientsByFoodID(foodID uint) ([]models.FoodNutrient, error) {
	return s.repo.GetByFoodID(foodID)
}

// GetFoodNutrientsByNutrientID retrieves food nutrients by nutrient ID
func (s *FoodNutrientService) GetFoodNutrientsByNutrientID(nutrientID uint) ([]models.FoodNutrient, error) {
	return s.repo.GetByNutrientID(nutrientID)
}

// UpdateFoodNutrient updates a food nutrient record
func (s *FoodNutrientService) UpdateFoodNutrient(foodNutrient *models.FoodNutrient) error {
	return s.repo.Update(foodNutrient)
}

// DeleteFoodNutrient removes a food nutrient record
func (s *FoodNutrientService) DeleteFoodNutrient(id uint) error {
	return s.repo.Delete(id)
} 