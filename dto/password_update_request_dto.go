package dto

// UpdatePasswordRequestDTO represents the user password update data
type UpdatePasswordRequestDTO struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// AdminUpdatePasswordRequestDTO represents the admin password update data
type AdminUpdatePasswordRequestDTO struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
} 