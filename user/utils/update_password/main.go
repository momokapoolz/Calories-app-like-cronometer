package main

import (
	"fmt"
	"log"
	"os"

	"github.com/momokapoolz/caloriesapp/user/database"
	"github.com/momokapoolz/caloriesapp/user/utils"
)

func main() {
	// Connect to the database
	database.ConnectDatabase()

	// The email and password to update
	email := "momoka@email.com"
	password := "securepassword"

	// Generate bcrypt hash
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Update the password hash in the database using the correct column name
	result := database.DB.Exec(`UPDATE "User" SET "PasswordHash" = ? WHERE email = ?`, hashedPassword, email)
	if result.Error != nil {
		log.Fatalf("Failed to update password: %v", result.Error)
	}

	rowsAffected := result.RowsAffected
	if rowsAffected == 0 {
		log.Fatalf("No user found with email: %s", email)
	}

	fmt.Printf("Successfully updated password for user: %s\n", email)
	fmt.Printf("New password hash: %s\n", hashedPassword)

	os.Exit(0)
}
