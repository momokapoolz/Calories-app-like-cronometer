package services

import (
	"github.com/momokapoolz/caloriesapp/nutrient/models"
	"github.com/momokapoolz/caloriesapp/nutrient/repository"
)

type NutrientService struct {
	repo *repository.NutrientRepository
}

func NewNutrientService(repo *repository.NutrientRepository) *NutrientService {
	return &NutrientService{repo: repo}
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
