package services

import (
	"time"

	foodRepo "github.com/momokapoolz/caloriesapp/food/repository"
	"github.com/momokapoolz/caloriesapp/meal_log/models"
	"github.com/momokapoolz/caloriesapp/meal_log/repository"
	mealLogItemsModels "github.com/momokapoolz/caloriesapp/meal_log_items/models"
	mealLogItemsRepo "github.com/momokapoolz/caloriesapp/meal_log_items/repository"
	mealLogItemsServices "github.com/momokapoolz/caloriesapp/meal_log_items/services"
)

// MealLogService handles business logic for meal log operations
type MealLogService struct {
	repo               *repository.MealLogRepository
	mealLogItemsRepo   *mealLogItemsRepo.MealLogItemRepository
	mealLogItemService *mealLogItemsServices.MealLogItemService
}

// NewMealLogService creates a new meal log service instance
func NewMealLogService(repo *repository.MealLogRepository, mealLogItemsRepo *mealLogItemsRepo.MealLogItemRepository, foodRepository *foodRepo.FoodRepository) *MealLogService {
	mealLogItemService := mealLogItemsServices.NewMealLogItemService(mealLogItemsRepo, foodRepository)
	return &MealLogService{
		repo:               repo,
		mealLogItemsRepo:   mealLogItemsRepo,
		mealLogItemService: mealLogItemService,
	}
}

// CreateMealLog creates a new meal log record
func (s *MealLogService) CreateMealLog(mealLog *models.MealLog) error {
	return s.repo.Create(mealLog)
}

// CreateMealLogItem creates a new meal log item with automatic quantity calculation
func (s *MealLogService) CreateMealLogItem(item *mealLogItemsModels.MealLogItem) error {
	return s.mealLogItemService.CreateMealLogItem(item)
}

// GetMealLogWithItemsByID retrieves a meal log with its items by ID
func (s *MealLogService) GetMealLogWithItemsByID(id uint) (*models.MealLogWithItems, error) {
	mealLog, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	items, err := s.mealLogItemsRepo.GetByMealLogID(id)
	if err != nil {
		return nil, err
	}

	return &models.MealLogWithItems{
		MealLog: *mealLog,
		Items:   items,
	}, nil
}

// GetMealLogByID retrieves a meal log record by ID
func (s *MealLogService) GetMealLogByID(id uint) (*models.MealLog, error) {
	return s.repo.GetByID(id)
}

// GetMealLogsByUserID retrieves all meal logs for a specific user
func (s *MealLogService) GetMealLogsByUserID(userID uint) ([]models.MealLog, error) {
	return s.repo.GetByUserID(userID)
}

// GetMealLogsByUserIDAndDate retrieves meal logs for a specific user and date
func (s *MealLogService) GetMealLogsByUserIDAndDate(userID uint, date time.Time) ([]models.MealLog, error) {
	return s.repo.GetByUserIDAndDate(userID, date)
}

// GetMealLogsByUserIDAndDateRange retrieves meal logs for a specific user within a date range
func (s *MealLogService) GetMealLogsByUserIDAndDateRange(userID uint, startDate, endDate time.Time) ([]models.MealLog, error) {
	return s.repo.GetByUserIDAndDateRange(userID, startDate, endDate)
}

// UpdateMealLog updates a meal log record
func (s *MealLogService) UpdateMealLog(mealLog *models.MealLog) error {
	return s.repo.Update(mealLog)
}

// DeleteMealLog removes a meal log record and all its items
func (s *MealLogService) DeleteMealLog(id uint) error {
	// First delete all meal log items
	if err := s.mealLogItemsRepo.DeleteByMealLogID(id); err != nil {
		return err
	}

	// Then delete the meal log
	return s.repo.Delete(id)
}

// VerifyMealLogOwnership checks if the specified meal log belongs to the user
func (s *MealLogService) VerifyMealLogOwnership(mealLogID, userID uint) error {
	return s.mealLogItemService.VerifyMealLogOwnership(mealLogID, userID)
}
