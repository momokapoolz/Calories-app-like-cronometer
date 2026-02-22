package auth

// Claims represents the authenticated user's identity extracted from a JWT.
// This is the canonical claims struct used by AuthMiddleware and set into the Gin context.
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}
