package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
	jsoniter "github.com/json-iterator/go"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	dblogger "github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/logger"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/models"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
)

type MovieController struct {
	db       *database.MongoDB
	dbLogger *dblogger.DBLogger
}

func NewMovieController(db *database.MongoDB, dbLogger *dblogger.DBLogger) *MovieController {
	return &MovieController{
		db:       db,
		dbLogger: dbLogger,
	}
}

var geminiClient *genai.Client
var clientOnce sync.Once

// High-performance JSON but gin already optimizes the json, so unnecessary (used just for learning)
var json = jsoniter.ConfigFastest

var validate = validator.New()

// Helper for fast JSON marshal
func mustMarshal(v any) []byte {
	data, _ := json.Marshal(v)
	return data
}

func (mc *MovieController) GetMovies(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	search := c.Query("search")

	filter := bson.M{}

	if search != "" {
        filter = bson.M{
            "title": bson.M{
                "$regex":   search, // slow for large data sets, instead create index on "text" for title
                "$options": "i", // for case-insensitive
            },
        }
    }

	var movies []models.Movie

	movieCollection := mc.db.Collection("movies")

	cursor, err := movieCollection.Find(ctx, filter)

	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Failed to fetch movies", gin.H{
			"endpoint": "/GetMovies",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}
	defer cursor.Close(ctx)

	// Streaming instead of cursor.All()
	for cursor.Next(ctx) {
		var movie models.Movie

		if err := cursor.Decode(&movie); err != nil {
			mc.dbLogger.Alerts("ERROR", "Decode error", gin.H{
				"endpoint": "/GetMovies",
				"error":    err.Error(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			return
		}

		movies = append(movies, movie)
	}

	if err := cursor.Err(); err != nil {
		mc.dbLogger.Alerts("ERROR", "Cursor error", gin.H{
			"endpoint": "/GetMovies",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	// c.JSON(http.StatusOK, movies)

	//mustMarshal is not efficient and unsafe as it silently ignores error
	c.Data(http.StatusOK, "application/json", mustMarshal(movies))
}

func (mc *MovieController) GetMovie(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var movieID = c.Param("imdb_id")
	if movieID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID required."})
		return
	}

	movieCollection := mc.db.Collection("movies")

	var movie models.Movie

	if err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie); err != nil {
		mc.dbLogger.Alerts("ERROR", "Decode error", gin.H{
			"endpoint": "/GetMovie",
			"movie_id": movieID,
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again!."})
		return
	}

	c.Data(http.StatusOK, "application/json", mustMarshal(movie))
}

func (mc *MovieController) AddMovie(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var movie models.Movie

	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	//validate is an instance of a validator which is used to validate the incoming struct against all the validations we have defined in the models
	if err := validate.Struct(movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	movieCollection := mc.db.Collection("movies")

	result, err := movieCollection.InsertOne(ctx, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't add the Movie, try again!"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (mc *MovieController) AdminReviewUpdate(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	role, err := utilities.GetRolefromContect(c)
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Error in getting role", gin.H{
			"endpoint": "/AdminReviewUpdate",
			"error":    err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be Admin to access this functionality"})
		return
	}

	movieId := c.Param("imdb_id")

	if movieId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID missing"})
		return
	}

	var req struct {
		AdminReview string `json:"admin_review"`
	}

	var resp struct {
		RankingName string `json:"ranking_name"`
		AdminReview string `json:"admin_review"`
	}

	if err := c.ShouldBind(&req); err != nil {
		mc.dbLogger.Alerts("ERROR", "Bind error", gin.H{
			"endpoint": "/AdminReviewUpdate",
			"movie_id": movieId,
			"error":    err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to bind for admin review"})
		return
	}

	sentiment, rankVal, err := mc.GetReviewRanking(req.AdminReview, c)
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Unable to get ranking for admin review", gin.H{
			"endpoint": "/AdminReviewUpdate",
			"movie_id": movieId,
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get ranking for admin review"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"admin_review": req.AdminReview,
			"ranking": bson.M{
				"ranking_value": rankVal,
				"ranking_name":  sentiment,
			},
		},
	}

	movieCollection := mc.db.Collection("movies")

	result, err := movieCollection.UpdateOne(ctx, bson.M{"imdb_id": movieId}, update)
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Error in Update query", gin.H{
			"endpoint": "/AdminReviewUpdate",
			"movie_id": movieId,
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update for admin review"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No movie found"})
		return
	}

	resp.AdminReview = req.AdminReview
	resp.RankingName = sentiment

	c.JSON(http.StatusOK, resp)
}

func (mc *MovieController) GetRecommendedMovies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	userId, err := utilities.GetUserIdfromContext(c)
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Error in getting userid", gin.H{
			"endpoint": "/GetRecommendedMovies",
			"error":    err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	favouriteGenres, err := mc.GetUserFavouriteGenres(userId, c)
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Error in getting favourite genres", gin.H{
			"endpoint": "/GetRecommendedMovies",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	var recommendedMovieLimit = os.Getenv("RECOMMENDED_MOVIE_LIMIT")

	var recommendedMovieLimitVal int64 = 5

	if recommendedMovieLimit != "" {
		recommendedMovieLimitVal, _ = strconv.ParseInt(recommendedMovieLimit, 10, 64)
	}

	findOptions := options.Find()

	findOptions.SetSort(bson.M{"ranking.ranking_value": 1})

	findOptions.SetLimit(recommendedMovieLimitVal)

	filter := bson.M{"genre.genre_name": bson.M{"$in": favouriteGenres}}

	movieCollection := mc.db.Collection("movies")

	cursor, err := movieCollection.Find(ctx, filter, findOptions)
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Error in executing query", gin.H{
			"endpoint": "/GetRecommendedMovies",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching recommended movies"})
		return
	}
	defer cursor.Close(ctx)

	var recommendedMovies []models.Movie

	for cursor.Next(ctx) {
		var recommendedMovie models.Movie

		if err := cursor.Decode(&recommendedMovie); err != nil {
			mc.dbLogger.Alerts("ERROR", "decoding query", gin.H{
				"endpoint": "/GetRecommendedMovies",
				"error":    err.Error(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error"})
			return
		}

		recommendedMovies = append(recommendedMovies, recommendedMovie)
	}

	if err := cursor.Err(); err != nil {
		mc.dbLogger.Alerts("ERROR", "cursor query", gin.H{
			"endpoint": "/GetRecommendedMovies",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(http.StatusOK, recommendedMovies)
}

func (mc *MovieController) GetUserFavouriteGenres(userId string, c *gin.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}

	opts := options.FindOne()

	opts.SetProjection(projection)

	var favoutiteGenres bson.M

	usercollection := mc.db.Collection("users")

	err := usercollection.FindOne(ctx, bson.M{"user_id": userId}, opts).Decode(&favoutiteGenres)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		} else {
			mc.dbLogger.Alerts("ERROR", "Error decoding or fetching favourite genres", gin.H{
				"endpoint": "/GetUserFavouriteGenres",
				"error":    err.Error(),
			})
			return []string{}, nil
		}
	}

	var favoutiteGenresNames []string

	genres, ok := favoutiteGenres["favourite_genres"].(primitive.A)
	if !ok {
		return []string{}, errors.New("Couldn't convert to array")
	}

	for _, genreNames := range genres {
		genreMap, ok := genreNames.(bson.M)
		if !ok {
			continue
		}

		if name, ok := genreMap["genre_name"].(string); ok {
			favoutiteGenresNames = append(favoutiteGenresNames, name)
		}
	}

	return favoutiteGenresNames, nil
}

func (mc *MovieController) GetReviewRanking(review string, c *gin.Context) (string, int, error) {
	rankings, err := mc.GetRankings(c)
	if err != nil {
		return "", 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // give Gemini more time
	defer cancel()

	var rankingNames []string

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			rankingNames = append(rankingNames, ranking.RankingName)
		}
	}

	allRankingNames := strings.Join(rankingNames, ",")

	var GEMINI_API_KEY = os.Getenv("GEMINI_API_KEY")
	var base_Prompt = os.Getenv("BASE_PROMPT")

	if GEMINI_API_KEY == "" {
		return "", 0, errors.New("could not find GEMINI SECRET KEY")
	}

	//Make this connection only once in main.go and pass it here to reduce time
	// client, err := genai.NewClient(ctx, option.WithAPIKey(GEMINI_API_KEY))

	client, err := GetGeminiClient(ctx)

	if err != nil {
		return "", 0, err
	}

	model := client.GenerativeModel("gemini-2.5-flash") //gemini-3-flash-preview

	basePrompt := strings.Replace(base_Prompt, "{rankings}", allRankingNames, 1)

	// response, err := llm.Call(ctx, basePrompt+review)
	resp, err := model.GenerateContent(ctx, genai.Text(basePrompt+review))

	if err != nil {
		return "", 0, err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", 0, errors.New("empty response from Gemini")
	}

	response := strings.TrimSpace(fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]))

	// Safety: normalize response
	response = strings.Trim(response, ".! \n\t")

	// logger.Info(response)

	rankVal := 0

	for _, r := range rankings {
		if strings.EqualFold(r.RankingName, response) {
			rankVal = r.RankingValue
			response = r.RankingName // normalize exact casing
			break
		}
	}

	// fallback (important
	if rankVal == 0 {
		response = "Okay"
		for _, r := range rankings {
			if r.RankingName == "Okay" {
				rankVal = r.RankingValue
				break
			}
		}
	}

	return response, rankVal, nil
}

func (mc *MovieController) GetRankings(c *gin.Context) ([]models.Ranking, error) {
	var rankings []models.Ranking

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	rankingCollection := mc.db.Collection("rankings")

	cursor, err := rankingCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, errors.New("Unable to fetch rankings")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ranking models.Ranking

		if err := cursor.Decode(&ranking); err != nil {
			return nil, errors.New("Error decoding ranking")
		}

		rankings = append(rankings, ranking)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.New("Cursor error")
	}

	return rankings, nil
}

func (mc *MovieController) GetGenres(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	genreCollection := mc.db.Collection("genres")

	var genreNames []models.Genre

	cursor, err := genreCollection.Find(ctx, bson.M{})
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Unable to fetch genres", gin.H{
			"endpoint": "/GetGenres",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var genreName models.Genre

		if err := cursor.Decode(&genreName); err != nil {
			mc.dbLogger.Alerts("WARN", "cursor decode error", gin.H{
				"endpoint": "/GetGenres",
				"error":    err.Error(),
			})
		}

		genreNames = append(genreNames, genreName)
	}

	if err := cursor.Err(); err != nil {
		mc.dbLogger.Alerts("WARN", "cursor error", gin.H{
			"endpoint": "/GetGenres",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(http.StatusOK, genreNames)
}

func GetGeminiClient(ctx context.Context) (*genai.Client, error) {
	var err error

	clientOnce.Do(func() {
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			err = errors.New("missing GEMINI_API_KEY")
			return
		}

		geminiClient, err = genai.NewClient(ctx, option.WithAPIKey(apiKey))
	})

	return geminiClient, err
}