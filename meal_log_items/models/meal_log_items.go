package models

// MealLogItem represents the meal_log_items table in the database
type MealLogItem struct {
	ID            uint    `gorm:"primaryKey;column:id" json:"id"`
	MealLogID     uint    `gorm:"column:meal_log_id;not null" json:"meal_log_id"`
	FoodID        uint    `gorm:"column:food_id;not null" json:"food_id"`
	Quantity      uint    `gorm:"column:quantity;not null" json:"quantity"`
	Quantity_     uint    `gorm:"column:quantity_;not null" json:"quantity_"` // Note: Field named as per schema
	QuantityGrams float64 `gorm:"column:quantity_grams;not null" json:"quantity_grams"`
}

// TableName specifies the table name for the MealLogItem model
func (MealLogItem) TableName() string {
	return "meal_log_items"
}
