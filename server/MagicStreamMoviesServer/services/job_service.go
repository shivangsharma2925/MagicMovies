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

func (js *JobService) UpdateStatus(imdbID string, status string, errMsg string) {

	filter := bson.M{"imdb_id": imdbID}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
			"error":      errMsg,
		},
	}

	js.jobCollection.UpdateOne(context.Background(), filter, update)
}

func (js *JobService) IncrementAttempts(imdbID string) int {

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
		context.Background(),
		filter,
		update,
		options,
	).Decode(&updatedJob)

	if err != nil {
		fmt.Printf("error in incrementing the attempts: %s", err)
	}

	return updatedJob.Attempts
}

func (js *JobService) MarkDone(imdbID string) {
	js.UpdateStatus(imdbID, "done", "")
}

func (js *JobService) MarkFailed(imdbID string, err string) {
	js.UpdateStatus(imdbID, "failed", err)
}
