package queue

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/tasks"
)

func EnqueueMovie(imdbID string) error {

	task, err := tasks.NewMovieTask(imdbID)

	if err != nil {
		return err
	}

	info, err := Client.Enqueue(
		task,

		// retry 5 times
		asynq.MaxRetry(5),

		// queue name
		asynq.Queue("movies"),
	)

	if err != nil {
		return err
	}

	log.Printf(
		"Movie task enqueued: taskID=%s imdbID=%s",
		info.ID,
		imdbID,
	)

	return nil
}
