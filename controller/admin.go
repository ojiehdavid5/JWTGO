package controller

import (
	"github.com/chuks/JWTGO/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/chuks/JWTGO/utils"
	// "log"
	"time"
	"fmt"
)

type Administrator struct {
	DB *gorm.DB
}

type Admin struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name		 string `gorm:"not null" json:"name"`
	Password string `gorm:"not null" json:"password,omitempty"`
}

func NewAdmin(db *gorm.DB) *Administrator {
	return &Administrator{DB: db}
}



func(a Administrator) Register(c *fiber.Ctx) error {
	var req Admin
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	admin := model.Admin{
		Name:        req.Name,
		PasswordHash: utils.GeneratePassword(req.Password),
	}
	fmt.Println(&admin)


	res := a.DB.Create(&admin)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "admin created",
		"admin":   admin,
	})

	
}
	



func (a Administrator) Login(c *fiber.Ctx) error {
	var req Admin
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	admin := model.Admin{}
	res := a.DB.Where("name = ?", req.Name).First(&admin)
	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}

	// Use req.Password instead of req.PasswordHash
	fmt.Println("Provided Password:", req.Password)
	fmt.Println("Stored Password Hash:", admin.PasswordHash)

	if !utils.VerifyPassword(admin.PasswordHash, req.Password) {
		return c.Status(401).JSON(fiber.Map{
			"message": "invalid password",
		})
	}

	token, err := utils.GenerateToken(admin.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"token":    token,
		"Login_at": time.Now(),
	})
}
func (a Administrator) GetUsers(c *fiber.Ctx) error {
	// Get all users from the database
	users := []model.User{}
	res := a.DB.Find(&users)
	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	// Return the users as JSON
	return c.Status(200).JSON(users)
}

func (a Administrator) DeleteUsers(c *fiber.Ctx) error{
	//GET THE BOOK ID  FROM URL PARAMETER
	userID := c.Params("id")
	//FIND THE BOOK FROM THE DATABASE
user:= model.User{}

res :=a.DB.First(&user,userID)

if res.Error != nil {
	return c.Status(404).JSON(fiber.Map{
		"message": res.Error.Error(),
	})
}

// delete user
res = a.DB.Delete(&user)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "user deleted",
	})
}






