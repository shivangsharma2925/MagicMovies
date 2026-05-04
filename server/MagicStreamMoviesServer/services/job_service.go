package services

import (
	"context"
	"time"

	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobService struct {
	db *database.MongoDB
	jobCollection *mongo.Collection
}

func NewJobService(db *database.MongoDB) *JobService {
	return &JobService{
		db: db,
		jobCollection: db.Collection("jobs"),
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

func (js *JobService) IncrementAttempts(imdbID string) {

	filter := bson.M{"imdb_id": imdbID}

	update := bson.M{
		"$inc": bson.M{
			"attempts": 1,
		},
	}

	js.jobCollection.UpdateOne(context.Background(), filter, update)
 
}

func (js *JobService) MarkDone(imdbID string) {
	js.UpdateStatus(imdbID, "done", "")
}

func (js *JobService) MarkFailed(imdbID string, err string) {
	js.UpdateStatus(imdbID, "failed", err)
}
