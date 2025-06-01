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
	// Find user first
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	// Generate bcrypt hash
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update user's password hash
	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

// ValidateCurrentPassword validates if the provided current password is correct
func (s *PasswordService) ValidateCurrentPassword(email, currentPassword string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	return utils.ComparePasswords(user.PasswordHash, currentPassword)
}
