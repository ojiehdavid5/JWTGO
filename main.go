package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/chuks/JWTGO/database"
	"log"
	"github.com/chuks/JWTGO/router"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	fmt.Println("DB_HOST:", os.Getenv("DB_HOST")) // Print environment variables
	fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))
	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))
	fmt.Println("DB_SSLMODE:", os.Getenv("DB_SSLMODE"))
	fmt.Println(time.Now())

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
 
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	router.SetupRoutes(app, db)

	log.Fatal(app.Listen(":3000"))
}
