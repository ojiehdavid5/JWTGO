package database

import (
	"fmt"
	"log"
	"os"

	"github.com/chuks/JWTGO/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect to the database
func Connect() (*gorm.DB, error) {
	// Construct the DSN string
	dsn := fmt.Sprintf(
		"host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)

	// Open the database connection
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err) // Log the detailed error
		return nil, fmt.Errorf("failed to connect to database: %w", err) // Return the error
	}

	fmt.Println("Connection Opened to Database")

	// AutoMigrate the schemas
	err = DB.AutoMigrate(&model.Book{}, &model.User{})
	if err != nil {
		log.Printf("Failed to migrate database: %v", err) // Log the detailed error
		return DB, fmt.Errorf("failed to migrate database: %w", err) // Return the error, but still return the DB connection (might be partially migrated)
	}

	fmt.Println("Database Migrated")
	return DB, nil // Return nil error if everything is successful
}