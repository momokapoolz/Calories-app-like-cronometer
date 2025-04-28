package main

import (
	"fmt"
	"github.com/momokapoolz/caloriesapp/user/database"
	"github.com/momokapoolz/caloriesapp/user/routes"
	"log"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>

	// Initialize database connection to PostgreSQL
	log.Println("Connecting to PostgreSQL database...")
	database.ConnectDatabase()

	// Set up API routes using Gin
	router := routes.SetupRoutes()

	// Start the Gin server
	port := "8080"
	log.Printf("Starting server on port %s...\n", port)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}

	fmt.Println("Calories App is running with PostgreSQL and Gin!")
}
