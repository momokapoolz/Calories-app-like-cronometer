package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	TokenKey    = "jwt_tokens"
)

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	URL      string
	Host     string
	Port     string
	Password string
	DB       int
}

// GetRedisConfig returns Redis configuration with sensible defaults
func GetRedisConfig() RedisConfig {

	url := os.Getenv("REDIS_URL")

	// Use local connection parameters
	//host := os.Getenv("REDIS_HOST")
	//if host == "" {
	//	host = "localhost"
	//}
	//
	//port := os.Getenv("REDIS_PORT")
	//if port == "" {
	//	port = "6379"
	//}
	//
	//password := os.Getenv("REDIS_PASSWORD")
	//
	//dbStr := os.Getenv("REDIS_DB")
	//db := 0
	//if dbStr != "" {
	//	if parsedDB, err := strconv.Atoi(dbStr); err == nil {
	//		db = parsedDB
	//	} else {
	//		log.Printf("Warning: Invalid REDIS_DB value '%s', using default DB 0", dbStr)
	//	}
	//}
	//
	//return RedisConfig{
	//	Host:     host,
	//	Port:     port,
	//	Password: password,
	//	DB:       db,
	//}

	//remote URL connection
	//if url != "" {
	//	opt, err := redis.ParseURL(url)
	//	if err != nil {
	//		log.Fatalf("Invalid REDIS_URL: %v", err)
	//	}
	//
	//	RedisClient = redis.NewClient(opt)
	//
	//	// Test the connection
	//	ctx := context.Background()
	//	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
	//		log.Fatalf("Failed to connect to Redis: %v", err)
	//	}
	//}
	return RedisConfig{URL: url}
}

// ConnectRedis initializes the Redis connection
func ConnectRedis() error {
	config := GetRedisConfig()

	// Use only local connection parameters
	//options := &redis.Options{
	//	Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
	//	Password: config.Password,
	//	DB:       config.DB,
	//}
	//
	//RedisClient = redis.NewClient(options)
	//
	//// Test the connection
	//ctx := context.Background()
	//_, err := RedisClient.Ping(ctx).Result()
	//if err != nil {
	//	return fmt.Errorf("failed to connect to Redis: %v", err)
	//}
	//
	//return nil

	var options *redis.Options
	var err error

	if config.URL != "" {
		options, err = redis.ParseURL(config.URL)
		if err != nil {
			return fmt.Errorf("invalid REDIS_URL: %v", err)
		}
	} else {
		options = &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
			Password: config.Password,
			DB:       config.DB,
		}
	}

	RedisClient = redis.NewClient(options)

	// Test the connection
	ctx := context.Background()
	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}

// StoreToken stores a JWT token in Redis and returns its ID
func StoreToken(token string, expiry time.Duration) (string, error) {
	ctx := context.Background()

	// Generate a random UUID for the token ID
	tokenID := uuid.New().String()

	// Store the token with the UUID as key
	key := fmt.Sprintf("%s:%s", TokenKey, tokenID)
	err := RedisClient.Set(ctx, key, token, expiry).Err()
	if err != nil {
		return "", err
	}

	return tokenID, nil
}

// GetToken retrieves a JWT token from Redis by ID
func GetToken(id string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", TokenKey, id)
	return RedisClient.Get(ctx, key).Result()
}

// DeleteToken removes a JWT token from Redis by ID
func DeleteToken(id string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", TokenKey, id)
	return RedisClient.Del(ctx, key).Err()
}
