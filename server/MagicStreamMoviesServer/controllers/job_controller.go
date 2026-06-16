package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	dblogger "github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/logger"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/models"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/queue"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genai"
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
			{Key: "_id", Value: 1},
			// Include other job fields you need here, for example:
			{Key: "imdb_id", Value: 1},
			{Key: "status", Value: 1},
			{Key: "attempts", Value: 1},
			{Key: "error", Value: 1},
			// Extract just the title from the joined movie document
			{Key: "title", Value: "$movie_details.title"},
		}}},
		bson.D{
			{Key: "$limit", Value: 30},
		},
	}

	// Execute the aggregation
	cursor, err := jobCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		jc.dbLogger.Alerts("Error", "Error in fetching jobs", gin.H{
			"enpoint": "/GetJobs",
			"error":   err.Error(),
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
			"error":   err.Error(),
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

func (jc *JobController) AskAiImdb(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	role, err := utilities.GetRolefromContect(c)
	if err != nil {
		jc.dbLogger.Alerts("ERROR", "Error in getting role", gin.H{
			"endpoint": "/AskAiImdb",
			"error":    err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong"})
		return
	}

	if role != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be Admin to access this functionality"})
		return
	}

	type AiResponse struct {
		MovieName string `json:"movie_name"`
		ImdbId    string `json:"imdbid"`
	}

	var AiImdbRequest struct {
		Prompt string `json:"prompt"`
	}

	var AiImdbResponse struct {
		Ids   []AiResponse `json:"ids,omitempty"`
		Error string       `json:"error,omitempty"`
	}

	if err := c.ShouldBind(&AiImdbRequest); err != nil {
		jc.dbLogger.Alerts("ERROR", "Error in binding request", gin.H{
			"endpoint": "/AskAiImdb",
			"error":    err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	if len(strings.TrimSpace(AiImdbRequest.Prompt)) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	client, err := GetGeminiClient(ctx)
	if err != nil {
		jc.dbLogger.Alerts("ERROR", "Error in connecting to AI model", gin.H{
			"endpoint": "/AskAiImdb",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect"})
		return
	}

	// model := client.GenerativeModel("gemini-2.5-flash")

	promptBytes, err := os.ReadFile("prompts/imdb_prompt.txt")
	if err != nil {
		jc.dbLogger.Alerts("ERROR", "Error in executing command", gin.H{
			"endpoint": "/AskAiImdb",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in executing command"})
		return
	}

	aiBasePrompt := string(promptBytes)

	fullPrompt := fmt.Sprintf("%s\n\nUser request: %s", aiBasePrompt, AiImdbRequest.Prompt)

	temp := float32(0)
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{GoogleSearch: &genai.GoogleSearch{}}, // Activates the grounding tool natively
		},
		Temperature: &temp, // Recommended by Google for optimal search grounding
	}

	rawText := ""
	modelArray := []string{"gemini-2.5-flash", "gemini-2.0-flash", "gemini-2.5-flash-lite", "gemini-3-flash-preview",}

	for _, model := range modelArray {

		jc.dbLogger.Log("INFO", "Trying Gemini model", gin.H{
			"endpoint": "/AskAiImdb",
			"model":    model,
		})

		resp, err := client.Models.GenerateContent(
			ctx,
			model,
			genai.Text(fullPrompt),
			config,
		)

		if err != nil {
			errText := strings.ToLower(err.Error())
			if 	strings.Contains(errText, "quota") ||
				strings.Contains(errText, "resource_exhausted") ||
				strings.Contains(errText, "429") {
				jc.dbLogger.Alerts("ERROR", "Quota exceeded", gin.H{
					"endpoint": "/AskAiImdb",
					"error":    err.Error(),
					"model":    model,
				})
				continue
			} else {
				jc.dbLogger.Alerts("ERROR", "Error in fetching response from AI", gin.H{
					"endpoint": "/AskAiImdb",
					"error":    err.Error(),
				})
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Request failed"})
				return
			}
		} else {
			rawText = strings.TrimSpace(resp.Text())
			break
		}
	}

	if rawText == "" {
		AiImdbResponse.Error = "no response from AI"
		c.JSON(http.StatusOK, AiImdbResponse)
		return
	}

	rawText = strings.TrimPrefix(rawText, "```json")
	rawText = strings.TrimPrefix(rawText, "```")
	rawText = strings.TrimSuffix(rawText, "```")
	rawText = strings.TrimSpace(rawText)

	var parsed struct {
		Invalid bool `json:"invalid"`
		Movies  []struct {
			MovieName string `json:"movie_name"`
			ImdbId    string `json:"imdbid"`
		} `json:"movies"`
	}

	if err := json.Unmarshal([]byte(rawText), &parsed); err != nil {
		jc.dbLogger.Alerts("ERROR", "Invalid JSON from AI", gin.H{
			"endpoint": "/AskAiImdb",
			"response": rawText,
			"error":    err.Error(),
		})
		AiImdbResponse.Error = "invalid AI response"
		c.JSON(http.StatusOK, AiImdbResponse)
		return
	}

	if parsed.Invalid {
		AiImdbResponse.Error = "invalid"
		c.JSON(http.StatusOK, AiImdbResponse)
		return
	}

	imdbPattern := regexp.MustCompile(`^tt\d+$`)
	seen := map[string]bool{}

	for _, movie := range parsed.Movies {
		movieName := strings.TrimSpace(movie.MovieName)
		imdbID := strings.TrimSpace(movie.ImdbId)

		if !imdbPattern.MatchString(imdbID) {
			continue
		}

		if seen[imdbID] {
			continue
		}

		seen[imdbID] = true

		AiImdbResponse.Ids = append(AiImdbResponse.Ids, AiResponse{
			MovieName: movieName,
			ImdbId:    imdbID,
		})
	}

	if len(AiImdbResponse.Ids) == 0 {
		AiImdbResponse.Error = "invalid"
		c.JSON(http.StatusOK, AiImdbResponse)
		return
	}

	c.JSON(http.StatusOK, AiImdbResponse)
}
