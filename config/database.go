package config

import (
	"fresh/app/models"
	"log"

	"gorm.io/gorm"
)

// InitDatabase initializes the database using the configuration system
func InitDatabase() *gorm.DB {
	// Load configuration
	config, err := LoadConfig("config/database.yml")
	if err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// Connect to database
	db, err := config.ConnectDatabase("")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema (this will add new columns but not drop existing data)
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Printf("Database migration failed: %v", err)
		log.Println("If you're getting constraint errors, you may need to manually fix the schema or reset the database with 'make db-reset'")
		log.Fatalf("Migration error: %v", err)
	}

	log.Println("Database migration completed")
	return db
}
