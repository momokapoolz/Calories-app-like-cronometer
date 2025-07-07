package dto

// LoginResponseDTO represents the response for both token-based and cookie-based auth
type LoginResponseDTO struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// TokenResponseDTO represents the response for token-based authentication
type TokenResponseDTO struct {
	User      interface{} `json:"user"`
	Token     string      `json:"token"`
	ExpiresIn int64       `json:"expires_in"`
}

// CookieResponseDTO represents the response for cookie-based authentication
type CookieResponseDTO struct {
	User      interface{} `json:"user"`
	ExpiresIn int64       `json:"expires_in"`
} 