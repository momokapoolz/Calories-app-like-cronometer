package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/database"
	"github.com/momokapoolz/caloriesapp/routes"
	user_database "github.com/momokapoolz/caloriesapp/user/database"
	user_routes "github.com/momokapoolz/caloriesapp/user/routes"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load .env file")
	}

	// Initialize database connection to PostgreSQL
	log.Println("Connecting to PostgreSQL database...")
	db := database.ConnectDatabase()
	user_database.ConnectDatabase() // For user module

	// Initialize Redis connection
	log.Println("Connecting to Redis...")
	if err := auth.ConnectRedis(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Set up API routes using Gin
	router := routes.SetupRoutes(db)
	
	// Set up User routes (from existing module)
	userRouter := router.Group("/api/v1")
	user_routes.SetupRoutes(userRouter)

	// Start the Gin server
	port := "8080"
	log.Printf("Starting server on port %s...\n", port)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}

	fmt.Println("Calories App is running with PostgreSQL, Redis, and Gin!")
}
