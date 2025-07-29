package dto

// RegisterDTO represents the user registration request data
type RegisterDTO struct {
	Name          string  `json:"name" binding:"required"`
	Email         string  `json:"email" binding:"required,email"`
	Password      string  `json:"password" binding:"required,min=6"`
	Age           int64   `json:"age" binding:"required,min=1"`
	Gender        string  `json:"gender" binding:"required"`
	Weight        float64 `json:"weight" binding:"required,gt=0"`
	Height        float64 `json:"height" binding:"required,gt=0"`
	Goal          string  `json:"goal" binding:"required"`
	ActivityLevel string  `json:"activity_level" binding:"required"`
}
