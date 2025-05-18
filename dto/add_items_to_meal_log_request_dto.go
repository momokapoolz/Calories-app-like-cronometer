package dto

type AddItemsToMealLogRequestDTO struct {
	Items []MealLogItemDTO `json:"items" binding:"required,dive"`
}
