package utils
import (
	"fmt"
	"os"
)




type User struct {
	Username string
	Email    string
}

// WriteUsersToFile writes a slice of users to a text file
func WriteUsersToFile(users []User, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, user := range users {
		_, err := fmt.Fprintf(file, "Username: %s, Email: %s\n", user.Username, user.Email)
		if err != nil {
			return err
		}
	}

	return nil
}