package database

import (
	"log"

	maindatabase "github.com/momokapoolz/caloriesapp/database"
	"gorm.io/gorm"
)

// DB is the shared GORM instance for the user module.
// It is set by ConnectDatabase() and points to the same connection
// pool as the main database package — no second connection is opened.
var DB *gorm.DB

// ConnectDatabase binds the user module's DB variable to the application-wide
// database connection. Must be called after database.ConnectDatabase().
func ConnectDatabase() {
	if maindatabase.DB == nil {
		log.Fatal("[user/database] Main database not initialized. Call database.ConnectDatabase() first.")
	}
	DB = maindatabase.DB
	log.Println("[user/database] Using shared database connection")
}
