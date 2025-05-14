package services

import (
	"github.com/momokapoolz/caloriesapp/food/models"
	"github.com/momokapoolz/caloriesapp/food/repository"
)

// FoodService handles business logic for food operations
type FoodService struct {
	repo *repository.FoodRepository
}

// NewFoodService creates a new food service instance
func NewFoodService(repo *repository.FoodRepository) *FoodService {
	return &FoodService{repo: repo}
}

// CreateFood creates a new food record
func (s *FoodService) CreateFood(food *models.Food) error {
	return s.repo.Create(food)
}

// GetFoodByID retrieves a food record by ID
func (s *FoodService) GetFoodByID(id uint) (*models.Food, error) {
	return s.repo.GetByID(id)
}

// GetAllFoods retrieves all food records
func (s *FoodService) GetAllFoods() ([]models.Food, error) {
	return s.repo.GetAll()
}

// UpdateFood updates a food record
func (s *FoodService) UpdateFood(food *models.Food) error {
	return s.repo.Update(food)
}

// DeleteFood removes a food record
func (s *FoodService) DeleteFood(id uint) error {
	return s.repo.Delete(id)
} 