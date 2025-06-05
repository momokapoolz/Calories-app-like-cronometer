package models

// Food represents the food table in the database
type Food struct {
	ID              uint    `gorm:"primaryKey;column:id" json:"id"`
	Name            string  `gorm:"column:name;not null" json:"name"`
	ServingSizeGram float64 `gorm:"column:serving_size_gram;not null" json:"serving_size_gram"`
	Source          string  `gorm:"column:source;not null" json:"source"`
	ImageURL        string  `gorm:"column:image_url" json:"image_url"`
}

// TableName specifies the table name for the Food model
func (Food) TableName() string {
	return "food"
}
