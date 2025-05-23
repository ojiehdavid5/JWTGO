package utils

import (
	"fmt"
	"os"
)

// User represents a registered user
type User struct {
	Email        string
	PasswordHash string
}

// WriteUserToFile writes a single user to a text file
func WriteUserToFile(user User, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Email: %s, Password: %s\n", user.Email, user.PasswordHash)
	return err
}