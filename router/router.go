package router

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/contrib/jwt"
	"gorm.io/gorm"
	"github.com/chuks/JWTGO/middleware"
)
import "github.com/chuks/JWTGO/controller"

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	api := app.Group("/api")
	book := controller.NewBook(db)
	auth := controller.NewAuth(db)
	// Book
	bookRoute := api.Group("/books")
	bookRoute.Get("/", book.GetBooks)
	bookRoute.Get("/:id",book.GetBook)
	bookRoute.Post("/",  book.CreateBook)
	bookRoute.Patch("/:id", middleware.JWTProtected, book.UpdateBook)
	bookRoute.Delete("/:id",middleware.JWTProtected, book.DeleteBook)

	// Auth
	authRoute := api.Group("/auth")
	authRoute.Post("/login", auth.Login)
	authRoute.Post("/register", auth.Register)
}
