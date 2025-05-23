package models

import (
	"time"
	mealLogItemsModels "github.com/momokapoolz/caloriesapp/meal_log_items/models"
)

// MealLog represents the meal_log table in the database
type MealLog struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID    uint      `gorm:"column:user_id;not null" json:"user_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	MealType  string    `gorm:"column:meal_type;not null" json:"meal_type"`
}

// MealLogWithItems represents a meal log with its associated items
type MealLogWithItems struct {
	MealLog MealLog                           `json:"meal_log"`
	Items   []mealLogItemsModels.MealLogItem `json:"items"`
}

// TableName specifies the table name for the MealLog model
func (MealLog) TableName() string {
	return "meal_log"
}
