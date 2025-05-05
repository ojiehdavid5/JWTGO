package controller

import (
	"fmt"
	"strings"

	"github.com/chuks/JWTGO/model"
	"github.com/chuks/JWTGO/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
	"github.com/resend/resend-go/v2"
	// "log"
)

type Auth struct {
	DB *gorm.DB
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	if strings.Contains(req.Email, "@gmail.com") {
	} else {
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
	fmt.Println(&user);
	apiKey := "re_eY82hvAL_89v9aqGgpVf8qusCDXrjESvD"

    client := resend.NewClient(apiKey)

    params := &resend.SendEmailRequest{
        From:    "onboarding@resend.dev",
        To:      []string{user.Email},
        Subject: "Hello World",
        Html:    "<p>Congrats on sending your <strong>first email</strong>!</p>",
		
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
	apiKey = "re_eY82hvAL_89v9aqGgpVf8qusCDXrjESvD"

    client = resend.NewClient(apiKey)

    params = &resend.SendEmailRequest{
        From:    "onboarding@resend.dev",
        To:      []string{"ojiehdavid5@gmail.com"},
        Subject: "Hello World",
        Html:    "<p>Congrats on sending your <strong>first email</strong>!</p>",
		
    }

    sent, err = client.Emails.Send(params)

	if err != nil {
		panic(err)
	}
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
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"token":       token,
		"Login in at": time.Now(),
	})
}
