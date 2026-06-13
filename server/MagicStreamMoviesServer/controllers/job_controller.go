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
	"go.mongodb.org/mongo-driver/mongo"
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

// func (jc *JobController) GetJobs(c *gin.Context) {
// 	jobCollection := jc.db.Collection("jobs")

// 	cursor, _ := jobCollection.Find(context.Background(), bson.M{})

	// var jobs []models.Job
// 	cursor.All(context.Background(), &jobs)

// 	c.JSON(http.StatusOK, jobs)
// }

func (jc *JobController) GetJobs(c *gin.Context) {
	jobCollection := jc.db.Collection("jobs")

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// 1. SORT FIRST (1 for ascending, -1 for descending)
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "created_at", Value: -1},
		}}},
		// 2. Join jobs with movies collection
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "movies"},
			{Key: "localField", Value: "imdb_id"},
			{Key: "foreignField", Value: "imdb_id"},
			{Key: "as", Value: "movie_details"},
		}}},
		// 3. Flatten the movie_details array (since lookup returns an array)
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$movie_details"},
			{Key: "preserveNullAndEmptyArrays", Value: true}, // Keeps job even if no movie matches
		}}},
		// 4. Project only the fields you want
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			// Include other job fields you need here, for example:
			{Key: "imdb_id", Value: 1},
			{Key: "status", Value: 1},
			{Key: "attempts", Value: 1},
			{Key: "error", Value: 1},
			// Extract just the title from the joined movie document
			{Key: "title", Value: "$movie_details.title"},
		}}},
	}

	// Execute the aggregation
	cursor, err := jobCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		jc.dbLogger.Alerts("Error", "Error in fetching jobs", gin.H{
			"enpoint": "/GetJobs",
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	defer cursor.Close(context.Background())

	// Decode results into a dynamic slice or updated struct
	var jobs []models.Job
	if err := cursor.All(context.Background(), &jobs); err != nil {
		jc.dbLogger.Alerts("Error", "Error in decoding jobs", gin.H{
			"enpoint": "/GetJobs",
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func (jc *JobController) RetryJob(c *gin.Context) {
	imdbID := c.Param("imdb_id")

	jobCollection := jc.db.Collection("jobs")

	jobCollection.UpdateOne(context.Background(),
		bson.M{"imdb_id": imdbID},
		bson.M{"$set": bson.M{"status": "pending"}},
	)

	queue.EnqueueMovie(imdbID)

	c.JSON(200, gin.H{"message": "Retry queued"})
}
