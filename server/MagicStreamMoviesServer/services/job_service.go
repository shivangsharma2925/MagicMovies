package services

import (
	"context"
	"fmt"
	"time"

	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JobService struct {
	db            *database.MongoDB
	jobCollection *mongo.Collection
}

func NewJobService(db *database.MongoDB) *JobService {

	JobService := &JobService{
		db:            db,
		jobCollection: db.Collection("jobs"),
	}

	JobService.UpdateStaleJobs()

	return JobService
}

func (js *JobService) UpdateStaleJobs() {
	filter := bson.M{
		"status": "processing",
	}

	update := bson.M{
		"$set": bson.M{
			"status": "failed",
		},
	}

	_, err := js.jobCollection.UpdateMany(context.Background(), filter, update)

	if err != nil {
		fmt.Printf("error in incrementing the attempts: %s", err)
		return
	}
}

func (js *JobService) UpdateStatus(imdbID string, status string, errMsg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"imdb_id": imdbID}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
			"error":      errMsg,
		},
	}

	result, err := js.jobCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update job %q to %q: %w", imdbID, status, err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("job %q not found while setting status to %q", imdbID, status)
	}
	return nil
}

func (js *JobService) IncrementAttempts(imdbID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"imdb_id": imdbID,
	}

	update := bson.M{
		"$inc": bson.M{
			"attempts": 1,
		},
	}

	options := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var updatedJob models.Job

	err := js.jobCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options,
	).Decode(&updatedJob)

	if err != nil {
		return 0, fmt.Errorf("increment attempts for job %q: %w", imdbID, err)
	}

	return updatedJob.Attempts, nil
}

func (js *JobService) MarkDone(imdbID string) error {
	return js.UpdateStatus(imdbID, string(models.StatusDone), "")
}

func (js *JobService) MarkFailed(imdbID string, err string) error {
	return js.UpdateStatus(imdbID, string(models.StatusFailed), err)
}

func (js *JobService) UpdateTitle(imdbId string, title string) {
	filter := bson.M{
		"imdb_id": imdbId,
	}

	update := bson.M{
		"$set": bson.M{
			"title": title,
		},
	}

	js.jobCollection.UpdateOne(context.Background(), filter, update)
}
