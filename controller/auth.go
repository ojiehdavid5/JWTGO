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

	// Validate email format
	if !strings.Contains(req.Email, "@gmail.com") && !strings.Contains(req.Email, "@yahoo.com") {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid email",
		})
	}

	// Check if the user already exists
	var existingUser model.User
	res := a.DB.Where("email = ?", req.Email).First(&existingUser)
	if res.Error == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "user already exists",
		})
	}

	// Create a new user
	user := model.User{
		Email:        req.Email,
		PasswordHash: utils.GeneratePassword(req.Password),
	}

	res = a.DB.Create(&user)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}

	// Send welcome email
	apiKey := os.Getenv("APIKEY")
	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{"ojiehdavid5@gmail.com"}, // Send to the registered user's email
		Subject: "Welcome to folben",
		Html:    "<h1> welcome</h1>",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
// Convert model.User to utils.User
	utilsUser := utils.User{
		Email:        req.Email,
		PasswordHash: user.PasswordHash, // Or however you want to handle password
	}



	if err := utils.WriteUserToFile(utilsUser, "user.txt"); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to write user to file",
		})
	}

	// Print the response
	fmt.Println("Email sent successfully:", sent)


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
		return c.Status(401).JSON(fiber.Map{
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
	Html:    fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>OTP Verification</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			background-color: #f4f4f4;
			margin: 0;
			padding: 0;
		}
		.container {
			max-width: 600px;
			margin: 0 auto;
			background: #ffffff;
			padding: 20px;
			border-radius: 5px;
			box-shadow: 0 2px 10px rgba(0,0,0,0.1);
		}
		h1 {
			color: #333333;
			text-align: center;
		}
		p {
			color: #555555;
			line-height: 1.5;
		}
		.otp {
			font-size: 24px;
			font-weight: bold;
			color: #007bff;
			text-align: center;
			margin: 20px 0;
		}
		.footer {
			text-align: center;
			margin-top: 20px;
			font-size: 12px;
			color: #888888;
		}
		@media (max-width: 600px) {
			.container {
				padding: 10px;
			}
			.otp {
				font-size: 20px;
			}
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>OTP Verification</h1>
		<p>Hi there,</p>
		<p>Thank you for registering! Please use the following One-Time Password (OTP) to complete your verification:</p>
		<div class="otp">%s</div>
		<p>This OTP is valid for the next 10 minutes. If you did not request this, please ignore this email.</p>
		<div class="footer">
			<p>&copy; 2023 Your Company Name. All</p>
		</div>
	</div>
</body>
</html>`, otp),
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Print the response
	println("Email sent successfully:", sent.Id)
	fmt.Println("OTP sent to email:", otp)



	//Returning the OTP directly to the client is a security risk. Remove this in production.
	return c.Status(200).JSON(fiber.Map{
		// "otp": otp, //This is for testing purposes only
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
