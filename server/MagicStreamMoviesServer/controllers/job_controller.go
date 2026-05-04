package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	dblogger "github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/logger"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/models"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/queue"
	"go.mongodb.org/mongo-driver/bson"
)

type JobController struct {
	db       *database.MongoDB
	dbLogger *dblogger.DBLogger
}

func NewJobController(db *database.MongoDB, dbLogger *dblogger.DBLogger) *JobController {
	return &JobController{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (jc *JobController) GetJobs(c *gin.Context) {
	jobCollection := jc.db.Collection("jobs")

	cursor, _ := jobCollection.Find(context.Background(), bson.M{})

	var jobs []models.Job
	cursor.All(context.Background(), &jobs)

	c.JSON(http.StatusOK, jobs)
}

func (jc *JobController) RetryJob(c *gin.Context) {
	imdbID := c.Param("id")

	jobCollection := jc.db.Collection("jobs")

	jobCollection.UpdateOne(context.Background(),
		bson.M{"imdb_id": imdbID},
		bson.M{"$set": bson.M{"status": "pending"}},
	)

	queue.PushToQueue(imdbID)

	c.JSON(200, gin.H{"message": "Retry queued"})
}
