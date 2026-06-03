package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/services"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/tasks"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/websocket"
)

type MoviesProcessor struct {
	movieService *services.MovieService
	jobService   *services.JobService
}

func (p *MoviesProcessor) processMovies(ctx context.Context, t *asynq.Task) error {
	var payload tasks.MoviePayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	imdbID := payload.ImdbID

	// So that any panic is treated as normal error and retry 
	defer func() {
		if r := recover(); r != nil {

			panicErr := fmt.Errorf("panic: %v", r)

			p.jobService.MarkFailed(imdbID, panicErr.Error())

			websocket.BroadcastJobUpdate(gin.H{
				"type":    "job_update",
				"imdb_id": imdbID,
				"status":  "failed",
			})

			p.movieService.DbLogger.Alerts("ERROR", "Error in processing movie", map[string]any{
				"endpoint": "/processMovies",
				"movieID":  imdbID,
				"error":    panicErr.Error(),
			})
		}
	}()

	// MARK PROCESSING
	attempts := p.jobService.IncrementAttempts(imdbID)

	p.jobService.UpdateStatus(
		imdbID,
		"processing",
		"",
	)

	websocket.BroadcastJobUpdate(gin.H{
		"type":     "job_update",
		"imdb_id":  imdbID,
		"status":   "processing",
		"attempts": attempts,
	})

	// PROCESS MOVIE
	err := p.movieService.ProcessMovie(imdbID, ctx)

	if err != nil {

		p.movieService.DbLogger.Alerts("ERROR", "Error in processing movie", map[string]any{
			"endpoint": "/processMovies",
			"movieID":  imdbID,
			"error":    err.Error(),
		})

		p.jobService.MarkFailed(imdbID, err.Error())

		websocket.BroadcastJobUpdate(gin.H{
			"type":    "job_update",
			"imdb_id": imdbID,
			"status":  "failed",
		})

		return err
	}

	// SUCCESS
	p.jobService.MarkDone(imdbID)

	websocket.BroadcastJobUpdate(gin.H{
		"type":     "job_update",
		"imdb_id":  imdbID,
		"status":   "done",
		"attempts": attempts,
	})

	websocket.BroadcastJobUpdate(gin.H{
		"type": "new_movie",
	})

	log.Printf("Movie processed: %s", imdbID)

	return nil

}
