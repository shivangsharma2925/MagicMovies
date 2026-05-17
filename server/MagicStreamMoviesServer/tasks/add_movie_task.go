package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const TypeMovieAdd = "movie:add"

type MoviePayload struct {
	ImdbID string `json:"imdb_id"`
}

func NewMovieTask(imdbID string) (*asynq.Task, error) {
	payload, err := json.Marshal(MoviePayload{
		ImdbID: imdbID,
	})

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeMovieAdd, payload), nil
}
