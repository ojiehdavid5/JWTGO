package controller

import("fmt"
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
	var body authRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	user := model.User{
		Email:        body.Email,
		PasswordHash: utils.GeneratePassword(body.Password),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create user",
		})
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}