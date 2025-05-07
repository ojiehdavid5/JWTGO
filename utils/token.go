package utils

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
	})

	// Convert the JWT_SECRET environment variable to a byte slice
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	t, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return t, nil
}

func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}
