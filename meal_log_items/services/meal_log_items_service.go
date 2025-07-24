package services

import (
	"errors"
	"fmt"

	foodRepo "github.com/momokapoolz/caloriesapp/food/repository"
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
	repo     *repository.MealLogItemRepository
	foodRepo *foodRepo.FoodRepository
}

// NewMealLogItemService creates a new meal log item service instance
func NewMealLogItemService(repo *repository.MealLogItemRepository, foodRepo *foodRepo.FoodRepository) *MealLogItemService {
	return &MealLogItemService{
		repo:     repo,
		foodRepo: foodRepo,
	}
}

// CreateMealLogItem creates a new meal log item with automatic quantity_grams calculation
func (s *MealLogItemService) CreateMealLogItem(item *models.MealLogItem) error {
	// Auto-calculate QuantityGrams based on food serving size and quantity
	if err := s.calculateQuantityGrams(item); err != nil {
		// Log error but don't fail - use provided value
		fmt.Printf("Warning: Could not auto-calculate quantity_grams: %v\n", err)
	}
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

// UpdateMealLogItem updates a meal log item with automatic quantity_grams calculation
func (s *MealLogItemService) UpdateMealLogItem(item *models.MealLogItem) error {
	// Auto-calculate QuantityGrams based on food serving size and quantity
	if err := s.calculateQuantityGrams(item); err != nil {
		// Log error but don't fail - use provided value
		fmt.Printf("Warning: Could not auto-calculate quantity_grams: %v\n", err)
	}
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

// AddItemsToMealLog adds multiple items to an existing meal log with automatic calculation
func (s *MealLogItemService) AddItemsToMealLog(mealLogID uint, items []models.MealLogItem) ([]models.MealLogItem, error) {
	// Validate that all items have the same meal log ID
	for i := range items {
		if items[i].MealLogID != mealLogID {
			return nil, errors.New("all items must have the same meal log ID")
		}

		// Auto-calculate QuantityGrams for each item
		if err := s.calculateQuantityGrams(&items[i]); err != nil {
			// Log error but don't fail - use provided value
			fmt.Printf("Warning: Could not auto-calculate quantity_grams for item %d: %v\n", i, err)
		}
	}

	// Call repository method to add items in a transaction
	return s.repo.CreateBatch(items)
}

// calculateQuantityGrams automatically calculates the quantity_grams based on food serving size
func (s *MealLogItemService) calculateQuantityGrams(item *models.MealLogItem) error {
	// Get food information to access serving size
	food, err := s.foodRepo.GetByID(item.FoodID)
	if err != nil {
		return fmt.Errorf("failed to get food info: %w", err)
	}

	// Calculate total grams: quantity * serving_size_gram
	calculatedGrams := float64(item.Quantity) * food.ServingSizeGram

	// Use calculated grams if:
	// 1. QuantityGrams is 0 (not provided)
	// 2. Or if the provided value seems incorrect (simple heuristic)
	if item.QuantityGrams == 0 || s.shouldRecalculateGrams(item, calculatedGrams) {
		item.QuantityGrams = calculatedGrams
		fmt.Printf("Auto-calculated QuantityGrams: %d × %.2f = %.2f grams\n",
			item.Quantity, food.ServingSizeGram, calculatedGrams)
	}

	return nil
}

// shouldRecalculateGrams determines if we should override the provided QuantityGrams
func (s *MealLogItemService) shouldRecalculateGrams(item *models.MealLogItem, calculatedGrams float64) bool {
	// If quantity > 1 but grams seems like single serving, recalculate
	// Allow 10% tolerance for measurement variations
	tolerance := calculatedGrams * 0.1
	return item.Quantity > 1 &&
		item.QuantityGrams < (calculatedGrams-tolerance)
}

// VerifyMealLogOwnership checks if the specified meal log belongs to the user
func (s *MealLogItemService) VerifyMealLogOwnership(mealLogID, userID uint) error {
	// query the meal_log table to check ownership
	// database query through our repository
	belongsToUser, err := s.repo.VerifyMealLogOwnership(mealLogID, userID)
	if err != nil {
		return fmt.Errorf("error verifying meal log ownership: %w", err)
	}

	if !belongsToUser {
		return ErrUnauthorizedAccess
	}

	return nil
}
