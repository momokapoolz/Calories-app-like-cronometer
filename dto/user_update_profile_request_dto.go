package dto

// UserUpdateProfileRequestDTO represents the user profile update data
type UserUpdateProfileRequestDTO struct {
	Name          *string  `json:"name,omitempty"`
	Email         *string  `json:"email,omitempty" binding:"omitempty,email"`
	Age           *int64   `json:"age,omitempty"`
	Gender        *string  `json:"gender,omitempty"`
	Weight        *float64 `json:"weight,omitempty"`
	Height        *float64 `json:"height,omitempty"`
	Goal          *string  `json:"goal,omitempty"`
	ActivityLevel *string  `json:"activity_level,omitempty"`
} 