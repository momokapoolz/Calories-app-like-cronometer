package dto

import "time"

// DashboardResponseDTO represents the data structure for the user dashboard
type DashboardResponseDTO struct {
	Date                string              `json:"date"`
	TotalCalories       float64             `json:"total_calories"`
	NumberOfMeals       int                 `json:"number_of_meals"`
	MealLogs            []MealLogSummaryDTO `json:"meal_logs"`
	TotalMacronutrients MacronutrientsDTO   `json:"total_macronutrients,omitempty"`
}

// MealLogSummaryDTO represents the summary of a meal log for the dashboard
type MealLogSummaryDTO struct {
	ID            uint                 `json:"id"`
	MealType      string               `json:"meal_type"`
	CreatedAt     time.Time            `json:"created_at"`
	TotalCalories float64              `json:"total_calories"`
	FoodItems     []FoodItemSummaryDTO `json:"food_items"`
}

// FoodItemSummaryDTO represents a summary of food item in a meal log
type FoodItemSummaryDTO struct {
	ID            uint    `json:"id"`
	FoodID        uint    `json:"food_id"`
	FoodName      string  `json:"food_name"`
	Quantity      uint    `json:"quantity"`
	QuantityGrams float64 `json:"quantity_grams"`
	Calories      float64 `json:"calories"`
}

// MacronutrientsDTO represents the macronutrient breakdown
type MacronutrientsDTO struct {
	Protein      float64 `json:"protein"`
	Carbohydrate float64 `json:"carbohydrate"`
	Fat          float64 `json:"fat"`
}
