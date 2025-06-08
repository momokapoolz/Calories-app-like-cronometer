package main

import (
	"testing"
	"time"

	"github.com/momokapoolz/caloriesapp/user/database"
	"github.com/momokapoolz/caloriesapp/user/models"
	"github.com/momokapoolz/caloriesapp/user/repository"
)

func TestMain(m *testing.M) {
	// Setup test database
	setup()
	// Run tests
	m.Run()
	// No need to explicitly clean up PostgreSQL database
}

func setup() {
	// Initialize test database
	database.ConnectDatabase()

	// Clean users table before testing
	database.DB.Exec("TRUNCATE TABLE users RESTART IDENTITY")
}

func TestUserCRUD(t *testing.T) {
	// Clean up before test
	database.DB.Exec("TRUNCATE TABLE users RESTART IDENTITY")

	// Create a test user
	user := models.User{
		Name:          "Test User",
		Email:         "test@example.com",
		PasswordHash:  "testpass",
		Age:           25,
		Gender:        "Female",
		Weight:        65.5,
		Height:        165.0,
		Goal:          "Muscle Gain",
		ActivityLevel: "High",
		CreatedAt:     time.Now(),
	}

	// Test repository
	userRepo := repository.NewUserRepository()

	// Test create
	err := userRepo.Create(&user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Fatal("User ID should not be zero after creation")
	}

	// Test find by ID
	foundUser, err := userRepo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to find user by ID: %v", err)
	}

	if foundUser.Email != user.Email {
		t.Fatalf("Expected email %s, got %s", user.Email, foundUser.Email)
	}

	// Test update
	user.Name = "Updated Name"
	err = userRepo.Update(&user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify update
	updatedUser, err := userRepo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to find updated user: %v", err)
	}

	if updatedUser.Name != "Updated Name" {
		t.Fatalf("Expected name 'Updated Name', got %s", updatedUser.Name)
	}

	// Test delete
	err = userRepo.Delete(user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify deletion
	_, err = userRepo.FindByID(user.ID)
	if err == nil {
		t.Fatal("User should be deleted")
	}
}
