package utils_test
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/chuks/JWTGO/utils"
	"os"
	"github.com/golang-jwt/jwt/v5"
)
func TestGenerateToken(t *testing.T) {
	// Test with a valid user ID
	id := uint(1)
	token, err := utils.GenerateToken(id)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	assert.NoError(t, err,"expected error to be nil")
	
	

	// Verify that the token is not empty
	if token == "" {
		t.Fatal("Generated token is empty")
	}
	assert.NotEmpty(t, token)

	// Verify that the token can be parsed
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}
	assert.NoError(t, err,"expected error to be nil")

	if !parsedToken.Valid {
		t.Fatal("Parsed token is not valid")
	}
}
func TestVerifyToken(t *testing.T) {
	// Test with a valid token
	id := uint(1)
	token, err := utils.GenerateToken(id)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	assert.NoError(t, err,"expected error to be nil")

	isValid, err := utils.VerifyToken(token)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}
	assert.NoError(t, err,"expected error to be nil")

	if !isValid {
		t.Fatal("Token is not valid")
	}

	// Test with an invalid token
	isValid, err = utils.VerifyToken("invalid_token")
	if err == nil {
		t.Fatal("Expected an error for invalid token, but got none")
	}
	assert.Error(t, err,"expected error to be not nil")
}
