package controller_test


import (
	// "os"
	"testing"
	// "fmt"
	// "os"
	// "github.com/chuks/JWTGO/database"

	"github.com/chuks/JWTGO/controller"
	"github.com/chuks/JWTGO/model"
	// "github.com/chuks/JWTGO/utils"
	"github.com/gofiber/fiber/v2"
	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/stretchr/testify/assert"
		"gorm.io/driver/postgres"
	"net/http"
	"net/http/httptest"
	"bytes"

	
)
type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
func setupTestDB(t * testing.T)*gorm.DB { 
	dsn := "host=localhost port=5432 user=user password=Qwerty dbname=postgres sslmode=disable"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	assert.NoError(t, err, "failed to connect to database")
	DB.AutoMigrate(&model.User{})
	assert.NoError(t, err, "failed to migrate database")

	return DB
}
func setupFiberApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	auth := controller.NewAuth(db)

	app.Post("/register", auth.Register)
	app.Post("/login", auth.Login)
	app.Post("/verify-otp", auth.VerifyOTP)

	return app
}
func TestRegister(t *testing.T) {
	db := setupTestDB(t)
	app := setupFiberApp(db)

	reqBody := `{"email":"yenum@gmail.com", "password":"securepassword"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// First registration attempt
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	// Attempt to register the same user again
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	respDuplicate, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, respDuplicate.StatusCode)
}

	
	




