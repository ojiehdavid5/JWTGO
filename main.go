package main
import(
	"fmt"
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	// "golang.org/x/crypto/bcrypt"
	// "gorm.io/gorm"
	"log"
	"github.com/chuks/JWTGO/database"
	"os"





)


func main() {

	err := godotenv.Load()
 if err != nil {
  log.Println("Error loading .env file")
 }
 
	database.Connect()
	app:=fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	   })
log.Fatal(app.Listen(":3000"))	  

	

	
}


func GenerateToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	 "user_id": id,
	})
   
	t, err := token.SignedString(os.Getenv("JWT_SECRET"))
	if err != nil {
	 return "", err
	}
   
	return t, nil
   }
