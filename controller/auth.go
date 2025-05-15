package controller

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chuks/JWTGO/model"
	"github.com/chuks/JWTGO/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/resend/resend-go/v2"
	"gorm.io/gorm"
)

type Auth struct {
	DB *gorm.DB
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type otpRequest struct {
	OTP string `json:"otp"` // Corrected field name
}

func NewAuth(db *gorm.DB) *Auth {
	return &Auth{DB: db}
}

func (a Auth) Register(c *fiber.Ctx) error {
	var req authRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	user := model.User{
		Email:        req.Email,
		PasswordHash: utils.GeneratePassword(req.Password),
	}

	// Check if the user already exists
	if !strings.Contains(req.Email, "@gmail.com") && !strings.Contains(req.Email, "@yahoo.com") {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid email",
		})
	}
	var existingUser model.User
	res := a.DB.Where("email = ?", req.Email).First(&existingUser)
	if res.Error == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "user already exists",
		})
	}
	fmt.Println(&user)
	apiKey := os.Getenv("APIKEY")
	fmt.Println("API Key:", apiKey)

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{"ojiehdavid5@gmail.com"},
		Subject: "Welcome to folben" + user.Email,
		Html:    "<p>Welcome to <strong>folben</strong> where we make travelling seamless for all user's</p>",
	}

	sent, err := client.Emails.Send(params)

	if err != nil {
		panic(err)
	}
	// Print the response
	println("Email sent successfully:", sent)
	// Create a new message

	// Create the user
	res = a.DB.Create(&user)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}

	// Create a new message
	// apiKey = os.Getenv("APIKEY")

	// client = resend.NewClient(apiKey)

	// params = &resend.SendEmailRequest{
	// 	From:    "onboarding@resend.dev",
	// 	To:      []string{"ojiehdavid5@gmail.com"},
	// 	Subject: "Hello World",
	// 	Html:    "<p>Congrats on sending your <strong>first email</strong>!</p>",
	// }
	// sent, err = client.Emails.Send(params)
	// if err != nil {
	// 	return c.Status(503).JSON(fiber.Map{
	// 		"message": err.Error(),
	// 	})
	// }
	// Print the response
	println("Email sent successfully:", sent)
	return c.Status(201).JSON(fiber.Map{
		"message": "user created",
	})
}

func (a Auth) Login(c *fiber.Ctx) error {
	var req authRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	var user model.User
	res := a.DB.Where("email = ?", req.Email).First(&user)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	if !utils.VerifyPassword(user.PasswordHash, req.Password) {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid password",
		})
	}

	otp, err := utils.SendOTP(user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	apiKey := os.Getenv("APIKEY")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{"ojiehdavid5@gmail.com"}, //This should be the user email
		Subject: " Your FOLBEN OTP IS " + otp,      //Concatenate the OTP to the subject
		Html:    "<p>This is just to key your account <strong>Safe</strong>!</p>",
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Print the response
	println("Email sent successfully:", sent)
	fmt.Println("OTP sent to email:", otp)

	//Returning the OTP directly to the client is a security risk. Remove this in production.
	return c.Status(200).JSON(fiber.Map{
		"otp": otp, //This is for testing purposes only
	})
}

func (a Auth) VerifyOTP(c *fiber.Ctx) error {
	var req otpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Get user email from request, you might want to change this
	email := c.Query("email")

	// Retrieve the user from the database using the email.
	var user model.User
	res := a.DB.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	valid, err := utils.VerifyOTP(user.Email, req.OTP) // Pass the actual OTP value
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if !valid {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid OTP",
		})
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	apiKey := os.Getenv("APIKEY")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{"ojiehdavid5@gmail.com"}, //This should be the user email
		Subject: "Welcome to folben",               //Concatenate the OTP to the subject
		Html:    "<p>Congrats on joining <strong>Folben</strong>!</p>",
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Print the response
	println("Email sent successfully:", sent)

	return c.Status(200).JSON(fiber.Map{
		"token":    token,
		"Login_at": time.Now(),
	})
}
