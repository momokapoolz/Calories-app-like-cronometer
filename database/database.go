package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/momokapoolz/caloriesapp/food/models"
	food_nutrients_models "github.com/momokapoolz/caloriesapp/food_nutrients/models"
	meal_log_models "github.com/momokapoolz/caloriesapp/meal_log/models"
	meal_log_items_models "github.com/momokapoolz/caloriesapp/meal_log_items/models"
	nutrient_models "github.com/momokapoolz/caloriesapp/nutrient/models"
	user_biometrics_models "github.com/momokapoolz/caloriesapp/user_biometrics/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() *gorm.DB {
	var err error

	// Use hardcoded DSN string
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

	// Open database connection
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		// Try to provide more detailed error message
		log.Printf("Error details: %v\n", err)
		log.Printf("Make sure PostgreSQL is running on port 5433 and the 'calorie_app_db' database exists")
		log.Fatal("Failed to connect to the database: ", err)
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

	// Handle existing tables if needed
	handleExistingTables()

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.Food{},
		&nutrient_models.Nutrient{},
		&food_nutrients_models.FoodNutrient{},
		&meal_log_models.MealLog{},
		&meal_log_items_models.MealLogItem{},
		&user_biometrics_models.UserBiometric{},
	)
	if err != nil {
		log.Fatal("Failed to auto migrate database schema: ", err)
	}

	fmt.Println("==================================")
	fmt.Println("Database connected and migrated successfully")
	fmt.Println("==================================")

	return DB
}

// handleExistingTables attempts to handle existing tables and their constraints
func handleExistingTables() {
	// List of tables to check
	tables := []string{"food", "nutrient", "food_nutrient", "meal_log", "meal_log_item", "user_biometric"}

	for _, table := range tables {
		// Check if table exists
		var count int64
		DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = ?", table).Count(&count)

		if count > 0 {
			log.Printf("Table %s exists, dropping related constraints...", table)

			// Get the list of constraints to drop
			var constraints []struct {
				ConstraintName string `gorm:"column:conname"`
				TableName      string `gorm:"column:relname"`
			}

			DB.Raw(`
				SELECT con.conname, cl.relname
				FROM pg_constraint con
				JOIN pg_class cl ON con.conrelid = cl.oid
				WHERE cl.relname = ?
			`, table).Scan(&constraints)

			// Drop each constraint
			for _, constraint := range constraints {
				log.Printf("Dropping constraint %s on table %s\n", constraint.ConstraintName, constraint.TableName)
				DB.Exec(fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s CASCADE", constraint.TableName, constraint.ConstraintName))
			}

			// Try to drop indexes if they exist
			log.Printf("Dropping indexes on %s table...", table)
			DB.Exec(fmt.Sprintf("DROP INDEX IF EXISTS idx_%s_id", table))
		}
	}
}
