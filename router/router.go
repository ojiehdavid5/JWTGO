package router

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/contrib/jwt"
)
import "github.com/chuks/JWTGO/controller"

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	// Book
	book := api.Group("/books")
	book.Get("/", controller.GetBooks)
	book.Get("/:id", controller.GetBook)
	book.Post("/", controller.CreateBook)
	book.Patch("/:id", controller.UpdateBook)
	book.Delete("/:id", controller.DeleteBook)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", controller.Login)
	auth.Post("/register", controller.Register)
}
