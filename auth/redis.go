package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	TokenKey    = "jwt_tokens"
)

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// GetRedisConfig returns Redis configuration with sensible defaults
func GetRedisConfig() RedisConfig {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	password := os.Getenv("REDIS_PASSWORD")
	if password == "" {
		log.Println("Warning: REDIS_PASSWORD not set in environment variables")
	}

	dbStr := os.Getenv("REDIS_DB")
	db := 0
	if dbStr != "" {
		if parsedDB, err := strconv.Atoi(dbStr); err == nil {
			db = parsedDB
		} else {
			log.Printf("Warning: Invalid REDIS_DB value '%s', using default DB 0", dbStr)
		}
	}

	return RedisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
	}
}

// ConnectRedis initializes the Redis connection
func ConnectRedis() error {
	config := GetRedisConfig()
	
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Test the connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}

// StoreToken stores a JWT token in Redis and returns its ID
func StoreToken(token string, expiry time.Duration) (int64, error) {
	ctx := context.Background()
	
	// Get the next available ID
	id, err := RedisClient.Incr(ctx, TokenKey+":counter").Result()
	if err != nil {
		return 0, err
	}

	// Store the token with the ID as key
	key := fmt.Sprintf("%s:%d", TokenKey, id)
	err = RedisClient.Set(ctx, key, token, expiry).Err()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetToken retrieves a JWT token from Redis by ID
func GetToken(id int64) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%d", TokenKey, id)
	return RedisClient.Get(ctx, key).Result()
}

// DeleteToken removes a JWT token from Redis by ID
func DeleteToken(id int64) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%d", TokenKey, id)
	return RedisClient.Del(ctx, key).Err()
} 