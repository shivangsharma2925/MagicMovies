package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/models"
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
		string(models.StatusProcessing),
		"",
	)

	websocket.BroadcastJobUpdate(gin.H{
		"type":     "job_update",
		"imdb_id":  imdbID,
		"status":   string(models.StatusProcessing),
		"attempts": attempts,
	})

	// PROCESS MOVIE
	
	err, title := p.movieService.ProcessMovie(imdbID, ctx)

	// UPDATE MOVIE NAME
	if len(title) != 0 && strings.Contains(imdbID, "tt"){
		go p.jobService.UpdateTitle(imdbID, title)
	}

	if err != nil {

		p.movieService.DbLogger.Alerts("ERROR", "Error in processing movie", map[string]any{
			"endpoint": "/processMovies",
			"movieID":  imdbID,
			"movieName": title,
			"error":    err.Error(),
		})

		p.jobService.MarkFailed(imdbID, err.Error())

		websocket.BroadcastJobUpdate(gin.H{
			"type":    "job_update",
			"imdb_id": imdbID,
			"title": title,
			"status":  string(models.StatusFailed),
		})

		return err
	}

	// SUCCESS
	p.jobService.MarkDone(imdbID)

	websocket.BroadcastJobUpdate(gin.H{
		"type":     "job_update",
		"imdb_id":  imdbID,
		"title": title,
		"status":   string(models.StatusDone),
		"attempts": attempts,
	})

	websocket.BroadcastJobUpdate(gin.H{
		"type": "new_movie",
	})

	log.Printf("Movie processed: %s", imdbID)
	p.movieService.DbLogger.Log("PROCESSED", "Movie processed correctly", map[string]any{
		"endpoint": "/processMovies",
		"movieID":  imdbID,
	})

	return nil

}
