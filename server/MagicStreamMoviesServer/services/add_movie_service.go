package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	dblogger "github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/logger"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/models"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MovieService struct {
	Db       *database.MongoDB
	DbLogger *dblogger.DBLogger
	genreMap map[string]models.Genre
}

func NewMovieService(db *database.MongoDB, logger *dblogger.DBLogger) *MovieService {
	ms := &MovieService{
		Db:       db,
		DbLogger: logger,
		genreMap: make(map[string]models.Genre),
	}

	err := ms.loadGenres()
	if err != nil {
		logger.Alerts("ERROR", "Failed to load genres", map[string]any{
			"error": err,
		})
		panic("Failed to Load Genres")
	}

	return ms
}

var HttpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func (ms *MovieService) loadGenres() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := ms.Db.Collection("genres")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	genreMap := make(map[string]models.Genre)

	for cursor.Next(ctx) {
		var g models.Genre
		if err := cursor.Decode(&g); err != nil {
			continue
		}

		genreMap[strings.ToLower(g.GenreName)] = g
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	ms.genreMap = genreMap
	return nil
}

func (ms *MovieService) ProcessMovie(imdbID string, ctx context.Context) (error, string) {

	// 1. Check duplicate
	exists, err := ms.MovieExists(imdbID, ctx)
	if err != nil {
		return err, ""
	}

	if exists {
		return nil, ""
	}

	// 2. Fetch data
	movieData, err := ms.FetchMovieData(imdbID, ctx)
	if err != nil {
		return err, ""
	}

	// 3. Transform
	movie := utilities.MapToSchema(movieData)

	// 4. Save
	return ms.SaveMovie(movie, ctx), movieData.Title
}

func (ms *MovieService) MovieExists(id string, ctx context.Context) (bool, error) {

	moviesCollection := ms.Db.Collection("movies")

	filter := bson.M{"imdb_id": id}
	if !strings.Contains(id, "tt") {
		id_int, _ := strconv.Atoi(id)
		filter = bson.M{
			"tmdb_id": id_int,
		}
	}

	err := moviesCollection.FindOne(ctx, filter).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		} else {
			ms.DbLogger.Alerts("ERROR", "Error finding if movie exists in database", map[string]any{
				"endpoint": "/MovieExists",
				"error":    err.Error(),
			})
			return false, err
		}
	}

	return true, nil
}

func (ms *MovieService) FetchMovieData(imdbID string, ctx context.Context) (*models.MovieData, error) {

	imdbIDnumeric, err := strconv.Atoi(imdbID)
	if err != nil {
		imdbIDnumeric = 0
	}

	tmdbID := imdbIDnumeric
	if strings.Contains(imdbID, "tt") {
		tmdbID, err = ms.getTMDBID(imdbID, ctx)
		if err != nil {
			// ms.dbLogger.Alerts("ERROR", "TMDB ID fetch failed", nil)
			return nil, err
		}
	}

	details, err := ms.getMovieDetails(tmdbID, ctx)
	if err != nil {
		return nil, err
	}

	youtubeID, trailerErr := ms.getTrailer(tmdbID, ctx) // optional failure

	if youtubeID == "" || trailerErr != nil {
		errMsg := ""
		if trailerErr != nil {
			errMsg = trailerErr.Error()
		}

		ms.DbLogger.Alerts("WARN", "No trailer Info found", map[string]any{
			"endpoint": "/FetchMovieData",
			"Imdb_ID":  details.ImdbID,
			"Tmdb_id":  tmdbID,
			"error":    errMsg,
		})
	}

	return &models.MovieData{
		ImdbID:       details.ImdbID,
		TmdbID:       tmdbID,
		Title:        details.Title,
		Poster:       "https://image.tmdb.org/t/p/w300" + details.PosterPath,
		Genres:       utilities.MapGenres(details.Genres, ms.genreMap),
		YoutubeID:    youtubeID,
		Ranking:      utilities.GetRanking(details.VoteAvg),
		IsAdult:      details.IsAdult,
		Overview:     details.Overview,
		ReleaseDate:  details.ReleaseDate,
		Tagline:      details.Tagline,
		Runtime:      details.Runtime,
		BackdropPath: "https://image.tmdb.org/t/p/w1280" + details.BackdropPath,
		VoteAvg:      details.VoteAvg,
	}, nil
}

func (ms *MovieService) getTMDBID(imdbID string, ctx context.Context) (int, error) {
	var tmdbID int

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/find/%s?external_source=imdb_id&api_key=%s",
		imdbID,
		os.Getenv("TMDB_SECRET_KEY"),
	)

	err := utilities.WithRetry(5, ctx, func() error {

		var result models.TMDBFindResponse

		err := utilities.DoRequest(ctx, url, &result, HttpClient)
		if err != nil {
			return err
		}

		if len(result.MovieResults) == 0 {
			ms.DbLogger.Alerts("WARN", "No Movie Results found", map[string]any{
				"endpoint": "/getTMDBID",
				"Imdb_ID":  imdbID,
			})
			return nil
		}

		tmdbID = result.MovieResults[0].ID
		return nil
	})

	return tmdbID, err
}

func (ms *MovieService) getMovieDetails(tmdbID int, ctx context.Context) (*models.TMDBMovieDetails, error) {
	var data models.TMDBMovieDetails

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/%d?api_key=%s",
		tmdbID,
		os.Getenv("TMDB_SECRET_KEY"),
	)

	err := utilities.WithRetry(3, ctx, func() error {

		return utilities.DoRequest(ctx, url, &data, HttpClient)
	})

	return &data, err
}

func (ms *MovieService) getTrailer(tmdbID int, ctx context.Context) (string, error) {
	var youtubeID string = ""

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/%d/videos?api_key=%s",
		tmdbID,
		os.Getenv("TMDB_SECRET_KEY"),
	)

	err := utilities.WithRetry(3, ctx, func() error {

		var data models.TMDBVideoResponse

		if err := utilities.DoRequest(ctx, url, &data, HttpClient); err != nil {
			return err
		}

		for _, v := range data.Results {
			if v.Site == "YouTube" && v.Type == "Trailer" {
				youtubeID = v.Key
				return nil
			}
		}

		return nil
	})

	return youtubeID, err
}

func (ms *MovieService) SaveMovie(movie *models.Movie, ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	collection := ms.Db.Collection("movies")

	_, err := collection.InsertOne(ctx, movie)
	if err != nil {

		if mongo.IsDuplicateKeyError(err) {
			return nil // safe ignore
		}

		ms.DbLogger.Alerts("ERROR", "Failed to insert movie", map[string]any{
			"imdb_id": movie.ImdbID,
			"error":   err.Error(),
		})

		return err
	}

	return nil
}

// Used when some bulk updates on movie collection needs to be done (called from setup_routes.go)
func (ms *MovieService) UpdateMovies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 50*time.Second)
	defer cancel()

	moviecollection := ms.Db.Collection("movies")

	opts := options.Find().SetProjection(bson.M{
		"_id":     0,
		"imdb_id": 1,
	})

	type moviesidstruct struct {
		ImdbId string `bson:"imdb_id"`
	}

	var moviesid []moviesidstruct

	cursor, err := moviecollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &moviesid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, id := range moviesid {
		err, _ := ms.ProcessMovie(id.ImdbId, ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "imdb_id": id.ImdbId})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": " all done"})

}
