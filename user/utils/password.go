package utils

import (
	"github.com/momokapoolz/caloriesapp/user/database"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePasswords compares a hashed password with a plain text password
func ComparePasswords(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// UpdateUserPassword updates a user's password with a new hashed password
func UpdateUserPassword(email string, newPassword string) error {
	// Generate bcrypt hash
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update the password in the database
	result := database.DB.Exec(`UPDATE "User" SET password = ? WHERE email = ?`, hashedPassword, email)
	return result.Error
}
