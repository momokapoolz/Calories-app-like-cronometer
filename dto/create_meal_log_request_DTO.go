package dto

type CreateMealLogRequestDTO struct {
	MealType string           `json:"meal_type"`
	Items    []MealLogItemDTO `json:"items"`
}

type MealLogItemDTO struct {
	FoodID        uint    `json:"food_id"`
	Quantity      uint    `json:"quantity"`
	QuantityGrams float64 `json:"quantity_grams"`
}
