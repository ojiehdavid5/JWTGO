package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/chuks/JWTGO/middleware"
)
import "github.com/chuks/JWTGO/controller"

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	api := app.Group("/api")
	book := controller.NewBook(db)
	auth := controller.NewAuth(db)
	admin := controller.NewAdmin(db)
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
	authRoute.Post("/otp", auth.VerifyOTP)


	// Admin
	adminRoute := api.Group("/admin")
	adminRoute.Post("/register", admin.Register)
	adminRoute.Post("/login", admin.Login)
	adminRoute.Get("/", middleware.JWTProtected, admin.GetUsers)
	adminRoute.Get("/:id",  admin.DeleteUsers)
}
