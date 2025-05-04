package auth

// Claims represents the JWT claims structure
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// LoginRequest represents the login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshTokenID int64 `json:"refresh_token_id" binding:"required"`
} 