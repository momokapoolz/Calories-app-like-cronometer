package services

import (
	"github.com/momokapoolz/caloriesapp/meal_log_items/models"
	"github.com/momokapoolz/caloriesapp/meal_log_items/repository"
)

// MealLogItemService handles business logic for meal log item operations
type MealLogItemService struct {
	repo *repository.MealLogItemRepository
}

// NewMealLogItemService creates a new meal log item service instance
func NewMealLogItemService(repo *repository.MealLogItemRepository) *MealLogItemService {
	return &MealLogItemService{repo: repo}
}

// CreateMealLogItem creates a new meal log item
func (s *MealLogItemService) CreateMealLogItem(item *models.MealLogItem) error {
	return s.repo.Create(item)
}

// GetMealLogItemByID retrieves a meal log item by ID
func (s *MealLogItemService) GetMealLogItemByID(id uint) (*models.MealLogItem, error) {
	return s.repo.GetByID(id)
}

// GetMealLogItemsByMealLogID retrieves all items for a specific meal log
func (s *MealLogItemService) GetMealLogItemsByMealLogID(mealLogID uint) ([]models.MealLogItem, error) {
	return s.repo.GetByMealLogID(mealLogID)
}

// GetMealLogItemsByFoodID retrieves all meal log items for a specific food
func (s *MealLogItemService) GetMealLogItemsByFoodID(foodID uint) ([]models.MealLogItem, error) {
	return s.repo.GetByFoodID(foodID)
}

// UpdateMealLogItem updates a meal log item
func (s *MealLogItemService) UpdateMealLogItem(item *models.MealLogItem) error {
	return s.repo.Update(item)
}

// DeleteMealLogItem removes a meal log item
func (s *MealLogItemService) DeleteMealLogItem(id uint) error {
	return s.repo.Delete(id)
}

// DeleteMealLogItemsByMealLogID removes all items for a specific meal log
func (s *MealLogItemService) DeleteMealLogItemsByMealLogID(mealLogID uint) error {
	return s.repo.DeleteByMealLogID(mealLogID)
} 