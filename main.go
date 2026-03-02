package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/joho/godotenv"

	_ "github.com/momokapoolz/caloriesapp/docs"

	"github.com/momokapoolz/caloriesapp/database"
	"github.com/momokapoolz/caloriesapp/routes"
	user_database "github.com/momokapoolz/caloriesapp/user/database"
	user_routes "github.com/momokapoolz/caloriesapp/user/routes"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           My API Swagger
// @version         1.0
// @description     API documentation
// @host            localhost:8080
// @BasePath        /api/v1
func main() {

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// INIT SENTRY
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Environment:      os.Getenv("APP_ENV"), // dev / staging / production
		TracesSampleRate: 1.0,                  // Performance monitoring
	})
	if err != nil {
		log.Fatalf("Sentry initialization failed: %v\n", err)
	}

	// Ensure buffered events are sent before shutdown
	defer sentry.Flush(2 * time.Second)

	log.Println("Connecting to PostgreSQL...")
	db := database.ConnectDatabase()
	user_database.ConnectDatabase()

	// Create router
	router := routes.SetupRoutes(db)

	// Attach Sentry middleware
	router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	apiV1 := router.Group("/api/v1")
	user_routes.SetupRoutes(apiV1)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := "8080"
	log.Printf("Starting server on :%s\n", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}

	fmt.Println("Calories App is running!")
}
