package utils_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/chuks/JWTGO/utils"
	"golang.org/x/crypto/bcrypt"

)
 

func TestGeneratePassword(t *testing.T) {
	password := "my_secure_password"
	hashedPassword := utils.GeneratePassword(password)

	// Verify that the hashed password is not empty
	assert.NotEmpty(t, hashedPassword)

	// Verify that the password can be verified
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err,"password verification failed")
}


func TestVerifyPassword(t *testing.T) {
	password := "my_secure_password"
	hashedPassword := utils.GeneratePassword(password)

	// Test successful verification
	assert.True(t, utils.VerifyPassword(hashedPassword, password),"password verification failed")
	// Test failed verification with an incorrect password
	assert.False(t, utils.VerifyPassword(hashedPassword, "wrong_password"),"password verification failed")
}
