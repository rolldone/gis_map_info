package redis_support

import (
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() *redis.Client {

	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisPort := os.Getenv("REDIS_PORT")
	redisDB := os.Getenv("REDIS_DB")

	// Convert string to int
	DB, _ := strconv.Atoi(redisDB)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprint(redisHost, ":", redisPort),
		Password: redisPassword,
		DB:       DB,
	})

	RedisClient = client
	return client
}
