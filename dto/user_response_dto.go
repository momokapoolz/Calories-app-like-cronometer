package dto

import "time"

// UserResponseDTO represents the user data returned in API responses
type UserResponseDTO struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Age           int64     `json:"age"`
	Gender        string    `json:"gender"`
	Weight        float64   `json:"weight"`
	Height        float64   `json:"height"`
	Goal          string    `json:"goal"`
	ActivityLevel string    `json:"activity_level"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
} 