package services

import (
	"github.com/momokapoolz/caloriesapp/user/repository"
	"github.com/momokapoolz/caloriesapp/user/utils"
)

type PasswordService struct {
	userRepo *repository.UserRepository
}

func NewPasswordService(userRepo *repository.UserRepository) *PasswordService {
	return &PasswordService{
		userRepo: userRepo,
	}
}

// UpdatePassword updates a user's password
func (s *PasswordService) UpdatePassword(email string, newPassword string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(user.ID, hashedPassword)
}

// ValidateCurrentPassword validates if the provided current password is correct
func (s *PasswordService) ValidateCurrentPassword(email, currentPassword string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	return utils.ComparePasswords(user.PasswordHash, currentPassword)
}
