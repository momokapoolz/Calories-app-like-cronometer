package dto

// LoginRequestDTO represents the user login data
type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email" example:"momomo@email.com"`
	Password string `json:"password" binding:"required,min=6" example:"momomo36"`
}
