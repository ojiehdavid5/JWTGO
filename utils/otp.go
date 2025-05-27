package utils

import (
	"crypto/rand"
	"context"
	"fmt"
	"github.com/chuks/JWTGO/database"
	"github.com/go-redis/redis/v8"
	"math/big"
	"time"
)


func GenerateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)
	for i := range otp {

		
		randInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[randInt.Int64()]
	}
	return string(otp), nil
}

// storeOTP stores the OTP in Redis with an expiration time.
func StoreOTP(ctx context.Context, userID, otp string, expiration time.Duration) error {

	key := fmt.Sprintf("otp:%s", userID) // Use a key specific to the user
	rdb, err := database.ConnectRedis()
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, otp, expiration).Err()
}
func verifyOTP(ctx context.Context, userID, otp string) (bool, error) {

	key := fmt.Sprintf("otp:%s", userID)
	rdb, err := database.ConnectRedis()
	if err != nil {
		return false, err
	}
	storedOTP, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		// Key does not exist (OTP has expired or is invalid)
		return false, nil
	} else if err != nil {
		return false, err
	}

	if storedOTP == otp {
		// OTP is valid, remove it from Redis to prevent reuse
		err := rdb.Del(ctx, key).Err()
		if err != nil {
			return true, err // Return true, but log the error
		}
		return true, nil
	}

	return false, nil
}
func SendOTP(userID string) (string, error) {
	// Generate a random OTP
	otp, err := GenerateOTP(6)
	if err != nil {
		return "", err
	}

	// Store the OTP in Redis with a 5-minute expiration time
	err = StoreOTP(context.Background(), userID, otp, 5*time.Minute)
	if err != nil {
		return "", err
	}

	// Send the OTP to the user (e.g., via email or SMS)
	fmt.Printf("Sending OTP %s to user %s\n", otp, userID)

	return otp, nil
}
func VerifyOTP(userID, otp string) (bool, error) {
	// Verify the OTP
	isValid, err := verifyOTP(context.Background(), userID, otp)
	if err != nil {
		return false, err
	}

	if isValid {
		fmt.Printf("OTP %s for user %s is valid\n", otp, userID)
	} else {
		fmt.Printf("OTP %s for user %s is invalid or expired\n", otp, userID)
	}

	return isValid, nil
}
