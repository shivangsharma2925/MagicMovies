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
		asynq.Queue(tasks.MovieQueueName()),
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

//function to keep track of data in worker queue(not called from anywhere)
func InspectMovieQueue(redisOpt asynq.RedisClientOpt) {
	inspector := asynq.NewInspector(redisOpt)
	defer inspector.Close()

	pending, pendingErr :=
		inspector.ListPendingTasks("movies")

	active, activeErr :=
		inspector.ListActiveTasks("movies")

	retry, retryErr :=
		inspector.ListRetryTasks("movies")

	archived, archivedErr :=
		inspector.ListArchivedTasks("movies")

	if pendingErr != nil {
		log.Printf("Inspect pending error: %v", pendingErr)
	}

	if activeErr != nil {
		log.Printf("Inspect active error: %v", activeErr)
	}

	if retryErr != nil {
		log.Printf("Inspect retry error: %v", retryErr)
	}

	if archivedErr != nil {
		log.Printf("Inspect archived error: %v", archivedErr)
	}

	log.Printf(
		"Movies queue: pending=%d active=%d retry=%d archived=%d",
		len(pending),
		len(active),
		len(retry),
		len(archived),
	)

	for _, task := range pending {
		log.Printf(
			"Pending task: id=%s type=%q queue=%q payload=%s",
			task.ID,
			task.Type,
			task.Queue,
			string(task.Payload),
		)
	}

	for _, task := range retry {
		log.Printf(
			"Retry task: id=%s type=%q error=%q nextProcessAt=%v",
			task.ID,
			task.Type,
			task.LastErr,
			task.NextProcessAt,
		)
	}

	for _, task := range archived {
		log.Printf(
			"Archived task: id=%s type=%q error=%q",
			task.ID,
			task.Type,
			task.LastErr,
		)
	}
}
