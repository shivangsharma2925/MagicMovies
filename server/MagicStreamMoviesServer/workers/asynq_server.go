package workers

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/services"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/tasks"
)

func StartWorkerServer(
	movieService *services.MovieService,
	jobService *services.JobService,
	redisOpt asynq.RedisClientOpt,
) {

	MoviesProcessor := &MoviesProcessor{
		movieService: movieService,
		jobService: jobService,
	}

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 3,

			Queues: map[string]int{
				"movies": 10,
			},
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeMovieAdd, MoviesProcessor.processMovies)

	log.Println("Asynq worker server started")

	if err := server.Run(mux); err != nil {
		log.Fatal(err)
	}
}
