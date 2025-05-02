package main

import (
	"fmt"
	"log"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Test DSN
	dsn := "host=localhost user=root password=1234 dbname=calories_app port=5433 sslmode=disable"
	
	fmt.Println("Trying to connect with DSN:", dsn)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	// Get SQL DB instance
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database connection:", err)
	}
	
	// Check connection
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	
	fmt.Println("Successfully connected to PostgreSQL!")
} 