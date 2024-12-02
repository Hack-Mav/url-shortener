package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

// RedisClient is an alias for the Redis client
type RedisClient = redis.Client

func ConnectCache() *redis.Client {
	// redisHost := os.Getenv("REDIS_HOST")
	redisHost := "redis-13327.c330.asia-south1-1.gce.redns.redis-cloud.com"
	// redisPort := os.Getenv("REDIS_PORT")
	redisPort := "13327"
	redisAddr := redisHost + ":" + redisPort

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Username: "default",
		Password: "gnbFTZJdKx0F33K7WPjqkVVafdYM6nKL",
	})

	_, err := client.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	return client
}
