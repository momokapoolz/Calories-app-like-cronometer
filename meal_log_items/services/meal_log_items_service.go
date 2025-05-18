package services

import (
	"errors"
	"fmt"

	"github.com/momokapoolz/caloriesapp/meal_log_items/models"
	"github.com/momokapoolz/caloriesapp/meal_log_items/repository"
)

// Error definitions
var (
	ErrMealLogNotFound    = errors.New("meal log not found")
	ErrUnauthorizedAccess = errors.New("unauthorized access to meal log")
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

// AddItemsToMealLog adds multiple items to an existing meal log
func (s *MealLogItemService) AddItemsToMealLog(mealLogID uint, items []models.MealLogItem) ([]models.MealLogItem, error) {
	// Validate that all items have the same meal log ID
	for i := range items {
		if items[i].MealLogID != mealLogID {
			return nil, errors.New("all items must have the same meal log ID")
		}
	}

	// Call repository method to add items in a transaction
	return s.repo.CreateBatch(items)
}

// VerifyMealLogOwnership checks if the specified meal log belongs to the user
func (s *MealLogItemService) VerifyMealLogOwnership(mealLogID, userID uint) error {
	// We need to query the meal_log table to check ownership
	// Since we don't have direct access to the meal_log repository here,
	// we'll use a database query through our repository
	belongsToUser, err := s.repo.VerifyMealLogOwnership(mealLogID, userID)
	if err != nil {
		return fmt.Errorf("error verifying meal log ownership: %w", err)
	}

	if !belongsToUser {
		return ErrUnauthorizedAccess
	}

	return nil
}
