package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/momokapoolz/caloriesapp/user/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() {
	var err error

	// Use the postgres database
	dsn := os.Getenv("POSTGRES_DB_CONNECTION_STRING")

	// Debug log
	log.Printf("Connecting with DSN: %s\n", dsn)

	// Set up logger for detailed SQL logs
	newLogger := logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Use singular table names
		},
		DisableForeignKeyConstraintWhenMigrating: true, // Disable foreign key constraint when migrating
	}

	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)

	if err != nil {
		// Try to provide more detailed error message
		log.Printf("Error details: %v\n", err)
		log.Printf("Make sure PostgreSQL is running on port 5433 and the 'postgres' database exists")
		log.Fatal("Failed to connect to database:", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database connection:", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Try to handle the existing users table
	handleExistingTable()

	// Auto migrate the models
	log.Println("Running auto-migration...")
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("==================================")
	fmt.Println("Database connected and migrated successfully")
	fmt.Println("==================================")
}

// handleExistingTable attempts to handle existing users table
func handleExistingTable() {
	// Check if users table exists
	var count int64
	DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = 'users'").Count(&count)

	if count > 0 {
		log.Println("Users table exists, dropping related constraints...")

		// Get the list of constraints to drop
		var constraints []struct {
			ConstraintName string `gorm:"column:conname"`
			TableName      string `gorm:"column:relname"`
		}

		DB.Raw(`
			SELECT con.conname, cl.relname
			FROM pg_constraint con
			JOIN pg_class cl ON con.conrelid = cl.oid
			WHERE cl.relname = 'users'
		`).Scan(&constraints)

		// Drop each constraint
		for _, constraint := range constraints {
			log.Printf("Dropping constraint %s on table %s\n", constraint.ConstraintName, constraint.TableName)
			DB.Exec(fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s CASCADE", constraint.TableName, constraint.ConstraintName))
		}

		// Try to drop indexes
		log.Println("Dropping indexes on users table...")
		DB.Exec("DROP INDEX IF EXISTS idx_users_email")
		DB.Exec("DROP INDEX IF EXISTS users_email_key")
		DB.Exec("DROP INDEX IF EXISTS idx_email")
	}
}
