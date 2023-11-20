package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// RedisClient ...
var RedisClient *redis.Client

// InitializeRedisConnection ...
func InitializeRedisConnection(redisURL string) {
	options, err := redis.ParseURL(redisURL)

	if err != nil {
		log.Fatal("Error parsing redis url")
		os.Exit(1)
	}

	client := redis.NewClient(options)

	pingResponse, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error pinging redis connection")
		os.Exit(1)
	}
	fmt.Println("REDIS PING RESPONSE: ", pingResponse)
	RedisClient = client
}
