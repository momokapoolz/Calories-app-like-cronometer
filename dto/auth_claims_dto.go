package dto

// ClaimsDTO represents the JWT claims structure
type ClaimsDTO struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// RefreshRequestDTO represents a token refresh request
type RefreshRequestDTO struct {
	RefreshTokenID string `json:"refresh_token_id" binding:"required"`
} 