package workers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/queue"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/services"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/websocket"
)

func StartWorkers(ctx context.Context, movieServices *services.MovieService, jobServices *services.JobService, n int) {
	for i := 0; i < n; i++ {
		go worker(ctx, movieServices, jobServices, i)
	}
}

func worker(ctx context.Context, movieServices *services.MovieService, jobServices *services.JobService, id int) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d shutting down\n", id)
			return

		case imdbID, ok := <-queue.AddMovieQueue:
			if !ok {
				log.Printf("Worker %d: queue closed\n", id)
				return
			}
			jobCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			// mark processing
			jobServices.IncrementAttempts(imdbID)
			jobServices.UpdateStatus(imdbID, "processing", "")
			websocket.BroadcastJobUpdate(gin.H{
				"imdb_id": imdbID,
				"status":  "processing",
			})

			err := processMovie(jobCtx, movieServices, imdbID)

			cancel()

			if err != nil {
				log.Printf("Worker %d failed: %s, err: %v\n", id, imdbID, err)
				jobServices.MarkFailed(imdbID, err.Error())
				websocket.BroadcastJobUpdate(gin.H{
					"imdb_id": imdbID,
					"status":  "failed",
				})
			} else {
				log.Printf("Worker %d success: %s\n", id, imdbID)
				jobServices.MarkDone(imdbID)
				websocket.BroadcastJobUpdate(gin.H{
					"imdb_id": imdbID,
					"status":  "done",
				})
			}
		}
	}
}

func processMovie(ctx context.Context, movieServices *services.MovieService, imdbID string) error {

	// added this so that retry won't delay the shutdown too long
	select {
	case <-ctx.Done():
		log.Printf("Cancelled Processing for %s", imdbID)
		return errors.New("Cancelled request")
	default:
	}

	err := movieServices.ProcessMovie(imdbID, ctx)
	if err == nil {
		return nil
	}

	movieServices.DbLogger.Alerts("ERROR", "Error in processing movie", map[string]any{
		"endpoint": "/processMovie",
		"movieID":  imdbID,
		"error":    err.Error(),
	})

	return err
}
