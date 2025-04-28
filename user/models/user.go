package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;unique;column:id"`
	Name          string    `gorm:"type:varchar(255);not null;column:name"`
	Email         string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_email;column:email"`
	PasswordHash  string    `gorm:"type:varchar(255);not null;column:password"`
	Age           int64     `gorm:"type:bigint;not null;column:age"`
	Gender        string    `gorm:"type:varchar(255);not null;column:gender"`
	Weight        float64   `gorm:"type:double precision;not null;column:weight"`
	Height        float64   `gorm:"type:double precision;not null;column:height"`
	Goal          string    `gorm:"type:varchar(255);not null;column:goal"`
	ActivityLevel string    `gorm:"type:varchar(255);not null;column:activity_level"`
	CreatedAt     time.Time `gorm:"type:timestamp with time zone;not null;column:created_at"`
	Role          string    `gorm:"type:varchar(255);not null;column:role"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}
