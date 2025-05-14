package models

// FoodNutrient represents the food_nutrients table in the database
type FoodNutrient struct {
	ID            uint    `gorm:"primaryKey;column:id" json:"id"`
	FoodID        uint    `gorm:"column:food_id;not null" json:"food_id"`
	NutrientID    uint    `gorm:"column:nutrient_id;not null" json:"nutrient_id"`
	AmountPer100g float64 `gorm:"column:amount_per_100g;not null" json:"amount_per_100g"`
}

// TableName specifies the table name for the FoodNutrient model
func (FoodNutrient) TableName() string {
	return "food_nutrients"
}
