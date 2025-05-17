package controller_test

import (
	// "os"
	"testing"
	"fmt"
	"time"
	// "os"
	// "github.com/chuks/JWTGO/database"

	"github.com/chuks/JWTGO/controller"
	"github.com/chuks/JWTGO/model"
	// "github.com/chuks/JWTGO/utils"
	"github.com/gofiber/fiber/v2"
	// "gorm.io/driver/sqlite"
	"bytes"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
)

func setupTestDB(t *testing.T) *gorm.DB {
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
func TestLogin(t *testing.T) {
	db := setupTestDB(t)
	app := setupFiberApp(db)

	// First, register a user
	registerReqBody := `{"email":"yenum@gmail.com", "password":"securepassword"}`
	registerReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(registerReqBody))
	registerReq.Header.Set("Content-Type", "application/json")
	_, err := app.Test(registerReq)
	assert.NoError(t, err)
	// Now, test the login
	loginReqBody := `{"email":"yenum@gmail.com", "password":"securepassword"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(loginReqBody))
	loginReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(loginReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	// Test login with incorrect password
	incorrectLoginReqBody := `{"email":"yenum@gmail.com", "password":"secrepassword"}`
	incorrectLoginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(incorrectLoginReqBody))
	incorrectLoginReq.Header.Set("Content-Type", "application/json")
	incorrectResp, err := app.Test(incorrectLoginReq)
	assert.NoError(t, err)
	assert.Equal(t, 401, incorrectResp.StatusCode)
}
	func VerifyOTP(t *testing.T) {
	db := setupTestDB(t)
	app := setupFiberApp(db)

	// First, register a user
	registerReqBody := `{"email":"yenum@gmail.com", "password":"securepassword"}`
	registerReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(registerReqBody))
	registerReq.Header.Set("Content-Type", "application/json")
	_, err := app.Test(registerReq)
	assert.NoError(t, err)

	// Now, test the login
	loginReqBody := `{"email":"yenum@gmail.com", "password":"securepassword"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(loginReqBody))
	loginReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(loginReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Simulate sending an OTP after login (you should implement this in your app)
	otp := "123456" // This should be generated and sent to the user in a real scenario

	// Test OTP verification
	otpReqBody := fmt.Sprintf(`{"otp":"%s"}`, otp)
	otpReq := httptest.NewRequest(http.MethodPost, "/verify-otp", bytes.NewBufferString(otpReqBody))
	otpReq.Header.Set("Content-Type", "application/json")
	otpResp, err := app.Test(otpReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, otpResp.StatusCode)

	// Test OTP verification with incorrect OTP
	incorrectOtpReqBody := `{"otp":"654321"}`
	incorrectOtpReq := httptest.NewRequest(http.MethodPost, "/verify-otp", bytes.NewBufferString(incorrectOtpReqBody))
	incorrectOtpReq.Header.Set("Content-Type", "application/json")
	incorrectOtpResp, err := app.Test(incorrectOtpReq)
	assert.NoError(t, err)
	assert.Equal(t, 401, incorrectOtpResp.StatusCode)

	// Test OTP verification with expired OTP
	// Simulate an expired OTP by waiting for the expiration time
	time.Sleep(6 * time.Minute) // Make sure your OTP expiration logic aligns with this duration
	expiredOtpReqBody := fmt.Sprintf(`{"otp":"%s"}`, otp) // Reusing the same OTP for testing purposes
	expiredOtpReq := httptest.NewRequest(http.MethodPost, "/verify-otp", bytes.NewBufferString(expiredOtpReqBody))
	expiredOtpReq.Header.Set("Content-Type", "application/json")
	expiredOtpResp, err := app.Test(expiredOtpReq)
	assert.NoError(t, err)
	assert.Equal(t, 401, expiredOtpResp.StatusCode)
}