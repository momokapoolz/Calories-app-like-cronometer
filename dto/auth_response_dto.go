package dto

// LoginResponseDTO is the standard response envelope for all auth endpoints
type LoginResponseDTO struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// AuthResponseDataDTO is the data payload inside a successful login/register/refresh response.
// Tokens are not included here — they are sent as HttpOnly cookies by the server.
type AuthResponseDataDTO struct {
	User      UserResponseDTO `json:"user"`
	ExpiresIn int64           `json:"expires_in"` // seconds until access token expires
}
