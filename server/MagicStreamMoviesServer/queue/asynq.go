package queue

import (
	"log"
	"os"
	"strconv"

	"github.com/hibiken/asynq"
)

func NewRedisConnOpt() asynq.RedisClientOpt {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		log.Fatal("Invalid REDIS_DB")
	}

	return asynq.RedisClientOpt{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}
}
