package queue

import "github.com/hibiken/asynq"

var Client *asynq.Client

func InitAsynqClient(redisOpt asynq.RedisClientOpt) {

	Client = asynq.NewClient(redisOpt)
}