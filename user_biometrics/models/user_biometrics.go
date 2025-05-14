package models

import (
	"time"
)

// UserBiometric represents the user_biometrics table in the database
type UserBiometric struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID    uint      `gorm:"column:user_id;not null" json:"user_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	Type      string    `gorm:"column:type;not null" json:"type"`
	Value     float64   `gorm:"column:value;not null" json:"value"`
	Unit      string    `gorm:"column:unit;not null" json:"unit"`
}

// TableName specifies the table name for the UserBiometric model
func (UserBiometric) TableName() string {
	return "user_biometrics"
}
