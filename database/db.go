package database

import (
 "fmt"
 "github.com/chuks/JWTGO/model"
 "os"
 "log"

 "gorm.io/gorm"
 "gorm.io/driver/postgres"

)

// DB instance
var DB *gorm.DB

// Connect to the database
func Connect() {
 dsn := fmt.Sprintf(
  "host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable",
  os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
 )

 DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
 if err != nil {
  log.Println("Failed to connect to database")
  fmt.Println("Error:", err)
 }
 fmt.Println("Connection Opened to Database")

 // Migrate the schemas
 DB.AutoMigrate(&model.Book{}, &model.User{})
 fmt.Println("Database Migrated")
}