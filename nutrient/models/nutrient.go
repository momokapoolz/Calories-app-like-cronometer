package models

// Nutrient represents the nutrient table in the database
type Nutrient struct {
	ID       uint   `gorm:"primaryKey;column:id" json:"id"`
	Name     string `gorm:"column:name;not null" json:"name"`
	Category string `gorm:"column:category;not null" json:"category"`
}

// TableName specifies the table name for the Nutrient model
func (Nutrient) TableName() string {
	return "nutrient"
}
