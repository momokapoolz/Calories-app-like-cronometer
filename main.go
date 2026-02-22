package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/momokapoolz/caloriesapp/database"
	"github.com/momokapoolz/caloriesapp/routes"
	user_database "github.com/momokapoolz/caloriesapp/user/database"
	user_routes "github.com/momokapoolz/caloriesapp/user/routes"
)

func main() {
	// Load .env file (optional — app works without it using env var defaults)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Open the single PostgreSQL connection and run all AutoMigrations.
	// Both the main module and the user module share this connection.
	log.Println("Connecting to PostgreSQL...")
	db := database.ConnectDatabase()

	// Bind the user module's DB variable to the shared connection.
	// No second connection is opened.
	user_database.ConnectDatabase()

	// Register application routes
	router := routes.SetupRoutes(db)

	apiV1 := router.Group("/api/v1")
	user_routes.SetupRoutes(apiV1)

	port := "8080"
	log.Printf("Starting server on :%s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}

	fmt.Println("Calories App is running!")
}
