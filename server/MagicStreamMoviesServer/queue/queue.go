package queue

import "github.com/bytedance/gopkg/util/logger"

var AddMovieQueue = make(chan string, 100)

//select helps in concurrency

func PushToQueue(id string) {
	select {
	case AddMovieQueue <- id:
		logger.Info("Movie with Id Added to queue:", id)
	default:
		logger.Info("queue if full, skipping movie with Id:", id)
	}
}
