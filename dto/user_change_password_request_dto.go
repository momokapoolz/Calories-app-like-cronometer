package dto

// UserChangePasswordRequestDTO represents the password change data
type UserChangePasswordRequestDTO struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
} 