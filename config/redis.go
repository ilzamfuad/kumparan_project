package config

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

func InitRedis() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	host := fmt.Sprintf("%s:%s", redisHost, redisPort)
	Rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: redisPassword,
		DB:       0,
	})

	return Rdb
}
