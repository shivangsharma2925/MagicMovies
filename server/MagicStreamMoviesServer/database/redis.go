package database

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() *redis.Client {

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		log.Fatal("Invalid REDIS_DB")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()

	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	logger.Info("Redis connected successfully")

	RedisClient = client

	return client
}
