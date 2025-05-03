package controller

import(
	// "fmt"
"github.com/chuks/JWTGO/database"
"github.com/chuks/JWTGO/model"
"github.com/gofiber/fiber/v2"
"github.com/chuks/JWTGO/utils"


)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
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
	res := database.DB.Create(&user)
	if res.Error != nil {
	 return c.Status(400).JSON(fiber.Map{
	  "message": res.Error.Error(),
	 })
	}
	return c.Status(201).JSON(fiber.Map{
	 "message": "user created",
	})
   }
   

   func Login(c *fiber.Ctx) error {
	var req authRequest
	if err := c.BodyParser(&req); err != nil {
	 return c.Status(400).JSON(fiber.Map{
	  "message": err.Error(),
	 })
	}
	var user model.User
	res := database.DB.Where("email = ?", req.Email).First(&user)
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
								 "token":token			})}