package services

import (
	"fmt"
	"time"

	"github.com/momokapoolz/caloriesapp/dto"
	"github.com/momokapoolz/caloriesapp/food/repository"
	foodNutrientsRepository "github.com/momokapoolz/caloriesapp/food_nutrients/repository"
	mealLogRepository "github.com/momokapoolz/caloriesapp/meal_log/repository"
	mealLogItemsRepository "github.com/momokapoolz/caloriesapp/meal_log_items/repository"
	nutrientRepository "github.com/momokapoolz/caloriesapp/nutrient/repository"
)

// Nutrient IDs for macronutrients (these should match your database)
const (
	CaloriesNutrientID     = 1 // Update with actual ID for calories
	ProteinNutrientID      = 2 // Update with actual ID for protein
	CarbohydrateNutrientID = 3 // Update with actual ID for carbohydrate
	FatNutrientID          = 4 // Update with actual ID for fat
)

// DashboardService handles business logic for dashboard operations
type DashboardService struct {
	mealLogRepo       *mealLogRepository.MealLogRepository
	mealLogItemsRepo  *mealLogItemsRepository.MealLogItemRepository
	foodRepo          *repository.FoodRepository
	nutrientRepo      *nutrientRepository.NutrientRepository
	foodNutrientsRepo *foodNutrientsRepository.FoodNutrientRepository
}

// NewDashboardService creates a new dashboard service instance
func NewDashboardService(
	mealLogRepo *mealLogRepository.MealLogRepository,
	mealLogItemsRepo *mealLogItemsRepository.MealLogItemRepository,
	foodRepo *repository.FoodRepository,
	nutrientRepo *nutrientRepository.NutrientRepository,
	foodNutrientsRepo *foodNutrientsRepository.FoodNutrientRepository,
) *DashboardService {
	return &DashboardService{
		mealLogRepo:       mealLogRepo,
		mealLogItemsRepo:  mealLogItemsRepo,
		foodRepo:          foodRepo,
		nutrientRepo:      nutrientRepo,
		foodNutrientsRepo: foodNutrientsRepo,
	}
}

// GetUserDashboard retrieves dashboard data for a user on a specific date
func (s *DashboardService) GetUserDashboard(userID uint, date time.Time) (*dto.DashboardResponseDTO, error) {
	// Get all meal logs for the user on the specified date
	mealLogs, err := s.mealLogRepo.GetByUserIDAndDate(userID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal logs: %w", err)
	}

	// Create the response DTO
	dashboard := &dto.DashboardResponseDTO{
		Date:                date.Format("2006-01-02"),
		NumberOfMeals:       len(mealLogs),
		MealLogs:            make([]dto.MealLogSummaryDTO, 0, len(mealLogs)),
		TotalMacronutrients: dto.MacronutrientsDTO{},
	}

	var totalCalories float64

	// Process each meal log
	for _, mealLog := range mealLogs {
		mealLogSummary := dto.MealLogSummaryDTO{
			ID:        mealLog.ID,
			MealType:  mealLog.MealType,
			CreatedAt: mealLog.CreatedAt,
			FoodItems: []dto.FoodItemSummaryDTO{},
		}

		// Get meal log items for this meal log
		items, err := s.mealLogItemsRepo.GetByMealLogID(mealLog.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get meal log items: %w", err)
		}

		var mealCalories float64

		// Process each meal log item
		for _, item := range items {
			// Get food information
			food, err := s.foodRepo.GetByID(item.FoodID)
			if err != nil {
				return nil, fmt.Errorf("failed to get food: %w", err)
			}

			// Calculate calories for this food item
			calories, err := s.calculateCalories(item.FoodID, item.QuantityGrams)
			if err != nil {
				return nil, fmt.Errorf("failed to calculate calories: %w", err)
			}

			// Create food item summary
			foodItem := dto.FoodItemSummaryDTO{
				ID:            item.ID,
				FoodID:        item.FoodID,
				FoodName:      food.Name,
				Quantity:      item.Quantity,
				QuantityGrams: item.QuantityGrams,
				Calories:      calories,
			}

			mealLogSummary.FoodItems = append(mealLogSummary.FoodItems, foodItem)
			mealCalories += calories
		}

		mealLogSummary.TotalCalories = mealCalories
		dashboard.MealLogs = append(dashboard.MealLogs, mealLogSummary)
		totalCalories += mealCalories

		// Calculate macronutrients (optional)
		protein, carbs, fat, err := s.calculateMacronutrients(mealLog.ID)
		if err == nil { // Only update if calculation succeeds
			dashboard.TotalMacronutrients.Protein += protein
			dashboard.TotalMacronutrients.Carbohydrate += carbs
			dashboard.TotalMacronutrients.Fat += fat
		}
	}

	dashboard.TotalCalories = totalCalories

	return dashboard, nil
}

// calculateCalories calculates the calories for a specific food item and quantity
func (s *DashboardService) calculateCalories(foodID uint, grams float64) (float64, error) {
	// Get calorie nutrient for this food
	foodNutrient, err := s.foodNutrientsRepo.GetByFoodIDAndNutrientID(foodID, CaloriesNutrientID)
	if err != nil {
		// If no specific calorie data is found, use a basic estimate (4 calories per gram)
		return grams * 4, nil
	}

	// Calculate calories based on the amount per 100g and the actual grams consumed
	calories := (foodNutrient.AmountPer100g / 100) * grams
	return calories, nil
}

// calculateMacronutrients calculates the macronutrients for a meal log
func (s *DashboardService) calculateMacronutrients(mealLogID uint) (protein, carbs, fat float64, err error) {
	// This is a stub implementation
	// In a real implementation, you would calculate these values from the food nutrients
	// For now, we'll return dummy values
	return 0, 0, 0, nil
}
