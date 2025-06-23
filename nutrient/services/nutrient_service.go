package services

import (
	"fmt"
	"time"

	"github.com/momokapoolz/caloriesapp/dto"
	foodRepo "github.com/momokapoolz/caloriesapp/food/repository"
	foodNutrientsRepo "github.com/momokapoolz/caloriesapp/food_nutrients/repository"
	mealLogRepo "github.com/momokapoolz/caloriesapp/meal_log/repository"
	mealLogItemsRepo "github.com/momokapoolz/caloriesapp/meal_log_items/repository"
	"github.com/momokapoolz/caloriesapp/nutrient/models"
	"github.com/momokapoolz/caloriesapp/nutrient/repository"
)

const (
	EnergyNutrientID       = 1
	ProteinNutrientID      = 2
	FatNutrientID          = 3
	CarbohydrateNutrientID = 4
	FiberNutrientID        = 5
	CholesterolNutrientID  = 6
	Vitamin_A_NutrientID   = 7
	Vitamin_B12_NutrientID = 8
	CalciumNutrientID      = 9
	IronNutrientID         = 10
)

type NutrientService struct {
	repo              *repository.NutrientRepository
	mealLogRepo       *mealLogRepo.MealLogRepository
	mealLogItemsRepo  *mealLogItemsRepo.MealLogItemRepository
	foodRepo          *foodRepo.FoodRepository
	foodNutrientsRepo *foodNutrientsRepo.FoodNutrientRepository
}

func NewNutrientService(
	repo *repository.NutrientRepository,
	mealLogRepo *mealLogRepo.MealLogRepository,
	mealLogItemsRepo *mealLogItemsRepo.MealLogItemRepository,
	foodRepo *foodRepo.FoodRepository,
	foodNutrientsRepo *foodNutrientsRepo.FoodNutrientRepository,
) *NutrientService {
	return &NutrientService{
		repo:              repo,
		mealLogRepo:       mealLogRepo,
		mealLogItemsRepo:  mealLogItemsRepo,
		foodRepo:          foodRepo,
		foodNutrientsRepo: foodNutrientsRepo,
	}
}

func (s *NutrientService) CreateNutrient(nutrient *models.Nutrient) error {
	return s.repo.Create(nutrient)
}

func (s *NutrientService) GetNutrientByID(id uint) (*models.Nutrient, error) {
	return s.repo.GetByID(id)
}

func (s *NutrientService) GetAllNutrients() ([]models.Nutrient, error) {
	return s.repo.GetAll()
}

func (s *NutrientService) GetNutrientsByCategory(category string) ([]models.Nutrient, error) {
	return s.repo.GetByCategory(category)
}

func (s *NutrientService) UpdateNutrient(nutrient *models.Nutrient) error {
	return s.repo.Update(nutrient)
}

func (s *NutrientService) DeleteNutrient(id uint) error {
	return s.repo.Delete(id)
}

// CalculateUserNutritionByDateRange calculates nutrition for a user within a date range
func (s *NutrientService) CalculateUserNutritionByDateRange(userID uint, startDate, endDate time.Time) (*dto.NutritionSummaryDTO, error) {
	mealLogs, err := s.mealLogRepo.GetByUserIDAndDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal logs: %w", err)
	}

	summary := &dto.NutritionSummaryDTO{
		UserID:                 userID,
		DateRange:              fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		MacroNutrientBreakDown: []dto.MacronutrientBreakdownDTO{},
		MicroNutrientBreakDown: []dto.MicronutrientDTO{},
		MealBreakdown:          []dto.MealNutritionDTO{},
	}

	totalNutrients := make(map[uint]float64) // nutrient_id -> total amount

	// Process each meal log
	for _, mealLog := range mealLogs {
		mealNutrition, foodCount, err := s.calculateMealNutrition(mealLog.ID)
		if err != nil {
			continue // Skip meals with calculation errors
		}

		// Add to meal breakdown
		summary.MealBreakdown = append(summary.MealBreakdown, dto.MealNutritionDTO{
			MealLogID:    mealLog.ID,
			MealType:     mealLog.MealType,
			Date:         mealLog.CreatedAt.Format("2006-01-02"),
			Calories:     mealNutrition[EnergyNutrientID],
			Protein:      mealNutrition[ProteinNutrientID],
			Carbohydrate: mealNutrition[CarbohydrateNutrientID],
			Fat:          mealNutrition[FatNutrientID],
			FoodCount:    foodCount,
		})

		// Add to total nutrients
		for nutrientID, amount := range mealNutrition {
			totalNutrients[nutrientID] += amount
		}
	}

	// Set total calories
	summary.TotalCalories = totalNutrients[EnergyNutrientID]

	// Create macro nutrient breakdown
	macroBreakdown := dto.MacronutrientBreakdownDTO{
		Energy:       totalNutrients[EnergyNutrientID],
		Protein:      totalNutrients[ProteinNutrientID],
		Fat:          totalNutrients[FatNutrientID],
		Carbohydrate: totalNutrients[CarbohydrateNutrientID],
		Fiber:        totalNutrients[FiberNutrientID],
		Cholesterol:  totalNutrients[CholesterolNutrientID],
		Vitamin_A:    totalNutrients[Vitamin_A_NutrientID],
		Vitamin_B12:  totalNutrients[Vitamin_B12_NutrientID],
		Calcium:      totalNutrients[CalciumNutrientID],
		Iron:         totalNutrients[IronNutrientID],
	}
	summary.MacroNutrientBreakDown = append(summary.MacroNutrientBreakDown, macroBreakdown)

	// Build micronutrient breakdown for other nutrients
	for nutrientID, amount := range totalNutrients {
		// Skip the nutrients already included in macro breakdown
		if s.isMacroNutrient(nutrientID) {
			continue
		}

		nutrient, err := s.repo.GetByID(nutrientID)
		if err != nil {
			continue
		}

		summary.MicroNutrientBreakDown = append(summary.MicroNutrientBreakDown, dto.MicronutrientDTO{
			NutrientID:   nutrientID,
			NutrientName: nutrient.Name,
			Amount:       amount,
			Unit:         "g", // You might want to store unit in nutrient table
		})
	}

	return summary, nil
}

// CalculateUserNutritionByDate calculates nutrition for a user on a specific date
func (s *NutrientService) CalculateUserNutritionByDate(userID uint, date time.Time) (*dto.NutritionSummaryDTO, error) {
	endDate := date.Add(24 * time.Hour)
	summary, err := s.CalculateUserNutritionByDateRange(userID, date, endDate)
	if err != nil {
		return nil, err
	}
	summary.DateRange = date.Format("2006-01-02")
	return summary, nil
}

// calculateMealNutrition calculates total nutrition for a specific meal
func (s *NutrientService) calculateMealNutrition(mealLogID uint) (map[uint]float64, int, error) {
	// Get all meal log items for this meal
	items, err := s.mealLogItemsRepo.GetByMealLogID(mealLogID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get meal log items: %w", err)
	}

	nutrition := make(map[uint]float64)
	foodCount := len(items)

	// Process each food item
	for _, item := range items {
		// Get all nutrients for this food
		foodNutrients, err := s.foodNutrientsRepo.GetByFoodID(item.FoodID)
		if err != nil {
			continue // Skip foods with no nutrient data
		}

		// Calculate nutrition based on quantity consumed
		for _, foodNutrient := range foodNutrients {
			// Convert amount per 100g to amount for consumed quantity
			actualAmount := (foodNutrient.AmountPer100g / 100.0) * item.QuantityGrams
			nutrition[foodNutrient.NutrientID] += actualAmount
		}
	}

	return nutrition, foodCount, nil
}

// isMacroNutrient checks if a nutrient ID is a macronutrient
func (s *NutrientService) isMacroNutrient(nutrientID uint) bool {
	return nutrientID == EnergyNutrientID ||
		nutrientID == ProteinNutrientID ||
		nutrientID == FatNutrientID ||
		nutrientID == CarbohydrateNutrientID ||
		nutrientID == FiberNutrientID ||
		nutrientID == CholesterolNutrientID ||
		nutrientID == Vitamin_A_NutrientID ||
		nutrientID == Vitamin_B12_NutrientID ||
		nutrientID == CalciumNutrientID ||
		nutrientID == IronNutrientID
}

// CalculateMealNutrition calculates nutrition for a specific meal log with user validation
func (s *NutrientService) CalculateMealNutrition(mealLogID uint, userID uint) (*dto.MealNutritionDetailDTO, error) {
	// First, verify that the meal log belongs to the user
	mealLog, err := s.mealLogRepo.GetByID(mealLogID)
	if err != nil {
		return nil, fmt.Errorf("meal log not found: %w", err)
	}

	if mealLog.UserID != userID {
		return nil, fmt.Errorf("access denied: meal log does not belong to user")
	}

	// Calculate nutrition for this meal
	mealNutrients, foodCount, err := s.calculateMealNutrition(mealLogID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate meal nutrition: %w", err)
	}

	// Build the detailed response
	response := &dto.MealNutritionDetailDTO{
		MealLogID:              mealLogID,
		UserID:                 userID,
		MealType:               mealLog.MealType,
		Date:                   mealLog.CreatedAt.Format("2006-01-02"),
		TotalCalories:          mealNutrients[EnergyNutrientID],
		FoodCount:              foodCount,
		MacroNutrientBreakDown: []dto.MacronutrientBreakdownDTO{},
		MicroNutrientBreakDown: []dto.MicronutrientDTO{},
	}

	// Create macro nutrient breakdown
	macroBreakdown := dto.MacronutrientBreakdownDTO{
		Energy:       mealNutrients[EnergyNutrientID],
		Protein:      mealNutrients[ProteinNutrientID],
		Fat:          mealNutrients[FatNutrientID],
		Carbohydrate: mealNutrients[CarbohydrateNutrientID],
		Fiber:        mealNutrients[FiberNutrientID],
		Cholesterol:  mealNutrients[CholesterolNutrientID],
		Vitamin_A:    mealNutrients[Vitamin_A_NutrientID],
		Vitamin_B12:  mealNutrients[Vitamin_B12_NutrientID],
		Calcium:      mealNutrients[CalciumNutrientID],
		Iron:         mealNutrients[IronNutrientID],
	}
	response.MacroNutrientBreakDown = append(response.MacroNutrientBreakDown, macroBreakdown)

	// Build micronutrient breakdown for other nutrients
	for nutrientID, amount := range mealNutrients {
		// Skip the nutrients already included in macro breakdown
		if s.isMacroNutrient(nutrientID) {
			continue
		}

		nutrient, err := s.repo.GetByID(nutrientID)
		if err != nil {
			continue
		}

		response.MicroNutrientBreakDown = append(response.MicroNutrientBreakDown, dto.MicronutrientDTO{
			NutrientID:   nutrientID,
			NutrientName: nutrient.Name,
			Amount:       amount,
			Unit:         "g", // You might want to store unit in nutrient table
		})
	}

	return response, nil
}
