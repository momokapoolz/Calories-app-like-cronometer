package services

import (
	"github.com/momokapoolz/caloriesapp/nutrient/models"
	"github.com/momokapoolz/caloriesapp/nutrient/repository"
)

// NutrientService handles business logic for nutrient operations
type NutrientService struct {
	repo *repository.NutrientRepository
}

// NewNutrientService creates a new nutrient service instance
func NewNutrientService(repo *repository.NutrientRepository) *NutrientService {
	return &NutrientService{repo: repo}
}

// CreateNutrient creates a new nutrient record
func (s *NutrientService) CreateNutrient(nutrient *models.Nutrient) error {
	return s.repo.Create(nutrient)
}

// GetNutrientByID retrieves a nutrient record by ID
func (s *NutrientService) GetNutrientByID(id uint) (*models.Nutrient, error) {
	return s.repo.GetByID(id)
}

// GetAllNutrients retrieves all nutrient records
func (s *NutrientService) GetAllNutrients() ([]models.Nutrient, error) {
	return s.repo.GetAll()
}

// GetNutrientsByCategory retrieves nutrients by category
func (s *NutrientService) GetNutrientsByCategory(category string) ([]models.Nutrient, error) {
	return s.repo.GetByCategory(category)
}

// UpdateNutrient updates a nutrient record
func (s *NutrientService) UpdateNutrient(nutrient *models.Nutrient) error {
	return s.repo.Update(nutrient)
}

// DeleteNutrient removes a nutrient record
func (s *NutrientService) DeleteNutrient(id uint) error {
	return s.repo.Delete(id)
} 