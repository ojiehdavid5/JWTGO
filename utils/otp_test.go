package utils_test
import ("testing"

	"github.com/stretchr/testify/assert"
	"github.com/chuks/JWTGO/utils"
	"context"
	"github.com/chuks/JWTGO/database"
	"time"
	"fmt"
)
func TestGenerateOTP(t *testing.T) {
	// Test with a valid length
	otp, err := utils.GenerateOTP(6)
	assert.NoError(t, err)
	assert.Len(t, otp, 6)


}
func TestStoreOTP(t *testing.T) {
	ctx := context.Background()
	userID := "test_user"
	otp := "123456"
	expiration := 5 * time.Minute

	// Store the OTP
	err := utils.StoreOTP(ctx, userID, otp, expiration)
	assert.NoError(t, err)

	// Retrieve the OTP from Redis
	rdb, err := database.ConnectRedis()
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	storedOTP, err := rdb.Get(ctx, fmt.Sprintf("otp:%s", userID)).Result()
	assert.NoError(t, err)
	assert.Equal(t, otp, storedOTP)
}
func TestVerifyOTP(t *testing.T) {
	ctx := context.Background()
	userID := "test_user"
	otp:= "123456"
	// Store the OTP in Redis
	err := utils.StoreOTP(ctx, userID, otp, 5*time.Minute)
	assert.NoError(t, err)

	// Verify the OTP
	isValid, err := utils.VerifyOTP(userID, otp)
	assert.NoError(t, err)
	assert.True(t, isValid)

	

}

