package controller

import (
	"fmt"
	"github.com/chuks/JWTGO/model"
	"github.com/gofiber/fiber/v2"
	// "log"
	"gorm.io/gorm"
)

type Book struct {
	DB *gorm.DB
}

func NewBook(db *gorm.DB) *Book {
	return &Book{DB: db}
}
	
func (b Book) GetBook(c *fiber.Ctx) error {
	// Get the book ID from the URL parameter
	bookID := c.Params("id")
	// Find the book in the database
	book := model.Book{}
	
	res := b.DB.First(&book, bookID)
	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	// Return the book as JSON
	return c.Status(200).JSON(book)
}
func (b Book) GetBooks(c *fiber.Ctx) error {
	// Get all books from the database
	books := []model.Book{}
	res := b.DB.Find(&books)
	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	// Return the books as JSON
	return c.Status(200).JSON(books)
}
func (b Book) CreateBook(c *fiber.Ctx) error {
	// log.Println("CreateBook controller called")

	// Parse the request body
	var req bookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Create a new book
	book := model.Book{
		Title:  req.Title,
		Author: req.Author,
	}

	res := b.DB.Create(&book)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	fmt.Println("Parsed request body:", req)

	// Return the created book as JSON
	return c.Status(201).JSON(book)
	// return c.SendString("CreateBook called")
}
func (b Book) UpdateBook(c *fiber.Ctx) error {
	// Get the book ID from the URL parameter
	bookID := c.Params("id")
	// Find the book in the database
	book := model.Book{}
	res := b.DB.First(&book, bookID)
	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	// Parse the request body
	var req bookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	book.Title = req.Title
	book.Author = req.Author
	res = b.DB.Save(&book)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	return c.Status(200).JSON(book)
}
func (b Book) DeleteBook(c *fiber.Ctx) error {
	// Get the book ID from the URL parameter
	bookID := c.Params("id")
	// Find the book in the database
	book := model.Book{}
	res := b.DB.First(&book, bookID)
	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	// Delete the book from the database
	res = b.DB.Delete(&book)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "book deleted",
	})
}

type bookRequest struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
}
