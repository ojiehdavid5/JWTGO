package database

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

func ConnectRedis() (*redis.Client,error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Redis connection options

	redisPassword := os.Getenv("REDIS_PASSWORD") // Get Redis password from environment variable

	ctx := context.Background()

	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: redisPassword,    // No password set
		DB:       0,                // Use default DB
	})
	// Ping the Redis server to check if it's available
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")
	return rdb,err
}
