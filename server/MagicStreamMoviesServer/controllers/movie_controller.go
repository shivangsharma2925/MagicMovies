package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	dblogger "github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/logger"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/models"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/queue"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/services"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/genai"
)

type HeroMovieCache struct {
	mu          sync.Mutex
	chachedData map[string][]models.Movie
	expiresAt   time.Time
}

type MovieController struct {
	db             *database.MongoDB
	dbLogger       *dblogger.DBLogger
	HeroMovieCache HeroMovieCache
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

	search = utilities.NormalizeText(search)

	if search != "" && len(search) > 2 {
		mc.GetSearchedMovies(c, ctx, search)
		return
	} else {

		// Instead of page and limit pagination, using Cursor based pagination,
		// if page or limit is not send in the URL, it will take the default value of 1 and 10 respectively.
		// page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		// if err != nil || page < 1 {
		// 	page = 1
		// }

		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit < 1 {
			limit = 10
		}

		cursorID := c.Query("cursor")

		// skip := (page - 1) * limit

		filter := bson.M{}

		// cursor pagination
		if cursorID != "" {

			objectID, err := primitive.ObjectIDFromHex(cursorID)

			if err == nil {
				filter["_id"] = bson.M{
					"$lt": objectID,
				}
			}
		}

		findOptions := options.Find().
			// 	SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "_id", Value: -1}}) // sort descending wise on _id since mongo db '_id' is based on time stamp(latest entries has larger _id value)

		movieCollection := mc.db.Collection("movies")

		// NO SEARCH

		movies := []models.Movie{}

		cursor, err := movieCollection.Find(ctx, filter, findOptions)

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

		var nextCursor any = nil

		if len(movies) > 0 {
			nextCursor = movies[len(movies)-1].ID.Hex()
		}

		// c.JSON(http.StatusOK, movies)

		c.JSON(http.StatusOK, gin.H{
			"movies":     movies,
			"nextCursor": nextCursor,
			"hasMore":    len(movies) == limit,
		})
	}

	//mustMarshal is not efficient and unsafe as it silently ignores error
	// c.Data(http.StatusOK, "application/json", mustMarshal(movies))
}

// using ATLAS search for giving search results
func (mc *MovieController) GetSearchedMovies(c *gin.Context, ctx context.Context, search string) {
	movieCollection := mc.db.Collection("movies")

	var movies []models.Movie

	pipeline := mongo.Pipeline{
		bson.D{
			bson.E{Key: "$search", Value: bson.D{
				bson.E{Key: "index", Value: "movies_search"},
				bson.E{Key: "compound", Value: bson.D{
					bson.E{Key: "should", Value: bson.A{

						// 1. Exact title match (highest boost)
						bson.D{
							bson.E{Key: "text", Value: bson.D{
								bson.E{Key: "query", Value: search},
								bson.E{Key: "path", Value: "title"},
								bson.E{Key: "score", Value: bson.D{
									bson.E{Key: "boost", Value: bson.D{
										bson.E{Key: "value", Value: 10},
									}},
								}},
							}},
						},

						// 2. Prefix / autocomplete match
						bson.D{
							bson.E{Key: "autocomplete", Value: bson.D{
								bson.E{Key: "query", Value: search},
								bson.E{Key: "path", Value: "title"},
								bson.E{Key: "tokenOrder", Value: "sequential"},
								bson.E{Key: "score", Value: bson.D{
									bson.E{Key: "boost", Value: bson.D{
										bson.E{Key: "value", Value: 5},
									}},
								}},
							}},
						},

						// 3. Fuzzy match (lowest boost)
						bson.D{
							bson.E{Key: "autocomplete", Value: bson.D{
								bson.E{Key: "query", Value: search},
								bson.E{Key: "path", Value: "title"},
								bson.E{Key: "tokenOrder", Value: "sequential"},
								bson.E{Key: "fuzzy", Value: bson.D{
									bson.E{Key: "maxEdits", Value: 2},
								}},
								bson.E{Key: "score", Value: bson.D{
									bson.E{Key: "boost", Value: bson.D{
										bson.E{Key: "value", Value: 2},
									}},
								}},
							}},
						},
					}},
					bson.E{Key: "minimumShouldMatch", Value: 1},
				}},
			}},
		},

		bson.D{
			bson.E{Key: "$project", Value: bson.D{
				bson.E{Key: "imdb_id", Value: 1},
				bson.E{Key: "title", Value: 1},
				bson.E{Key: "poster_path", Value: 1},
				bson.E{Key: "youtube_id", Value: 1},
				bson.E{Key: "genre", Value: 1},
				bson.E{Key: "ranking", Value: 1},
				bson.E{Key: "admin_review", Value: 1},

				// Useful for debugging ranking.
				bson.E{Key: "score", Value: bson.D{
					bson.E{Key: "$meta", Value: "searchScore"},
				}},
			}},
		},

		bson.D{
			bson.E{Key: "$limit", Value: 20},
		},
	}

	cursor, err := movieCollection.Aggregate(ctx, pipeline)

	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Error in searching movies", gin.H{
			"endpoint": "/GetSearchedMovies",
			"error":    err.Error(),
		})

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search movies",
		})
		return
	}

	if err := cursor.All(ctx, &movies); err != nil {
		mc.dbLogger.Alerts("ERROR", "Error in parsing movies", gin.H{
			"endpoint": "/GetSearchedMovies",
			"error":    err.Error(),
		})

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search movies",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
	})
}

// Manually giving search results - not using anywhere
func (mc *MovieController) GetSearchedMoviesManually(c *gin.Context, ctx context.Context, search string) {

	var finalMovies []models.Movie

	var scoredMovies []models.ScoredMovie //for sending movies with best match first

	movieMap := make(map[string]bool)

	// STAGE 1 — PREFIX REGEX SEARCH

	prefixFilter := bson.M{
		"title_normalized": bson.M{
			"$regex":   "^" + regexp.QuoteMeta(search),
			"$options": "i",
		},
	}

	movieCollection := mc.db.Collection("movies")

	prefixCursor, err := movieCollection.Find(c, prefixFilter, options.Find().SetLimit(10))

	if err == nil {

		var prefixMovies []models.Movie

		prefixCursor.All(context.Background(), &prefixMovies)

		for _, movie := range prefixMovies {

			if !movieMap[movie.ImdbID] {

				movieMap[movie.ImdbID] = true

				scoredMovies = append(scoredMovies, models.ScoredMovie{
					Movie: movie,
					Score: 100,
				})
			}
		}
	}

	// STAGE 2 — TOKEN SEARCH FALLBACK

	if len(scoredMovies) < 10 {

		searchTokens := utilities.GenerateSearchTokens(search)

		if len(searchTokens) > 0 {

			tokenFilter := bson.M{
				"search_tokens": bson.M{
					"$in": searchTokens,
				},
			}

			tokenCursor, err := movieCollection.Find(c, tokenFilter, options.Find().SetLimit(20))

			if err == nil {

				var tokenMovies []models.Movie

				tokenCursor.All(context.Background(), &tokenMovies)

				for _, movie := range tokenMovies {

					if !movieMap[movie.ImdbID] {

						movieMap[movie.ImdbID] = true

						score := utilities.CalculateTokenScore(
							searchTokens,
							movie.SearchTokens,
						)

						// medium priority
						finalScore := float64(score * 10)

						scoredMovies = append(
							scoredMovies,
							models.ScoredMovie{
								Movie: movie,
								Score: finalScore,
							},
						)
					}
				}
			}
		}
	}

	// STAGE 3 — FUZZY SEARCH FALLBACK
	if len(scoredMovies) < 10 {

		searchTokens := utilities.GenerateSearchTokens(search)

		fuzzyThreshold := 0.7

		if len(searchTokens) == 1 {
			fuzzyThreshold = 0.45
		}

		if len(searchTokens) > 0 {

			fuzzyCursor, err := movieCollection.Find(
				context.Background(),
				bson.M{},
				options.Find().
					SetLimit(200).
					SetProjection(bson.M{
						"title":            1,
						"title_normalized": 1,
						"search_tokens":    1,
						"poster_path":      1,
						"imdb_id":          1,
						"youtube_id":       1,
						"genre":            1,
						"ranking":          1,
						"admin_review":     1,
					}),
			)

			var fuzzyMovies []models.Movie

			if err == nil {
				fuzzyCursor.All(context.Background(), &fuzzyMovies)

				for _, movie := range fuzzyMovies {

					if !movieMap[movie.ImdbID] {

						score := utilities.CalculateFuzzyTokenScore(searchTokens, movie.SearchTokens)

						if score > fuzzyThreshold {
							movieMap[movie.ImdbID] = true

							scoredMovies = append(scoredMovies, models.ScoredMovie{
								Movie: movie,
								Score: score,
							})
						}
					}
				}
			}
		}
	}

	slices.SortStableFunc(scoredMovies, func(a, b models.ScoredMovie) int {
		if a.Score > b.Score {
			return -1
		}
		if a.Score < b.Score {
			return 1
		}
		return 0
	})

	for _, scoredMovie := range scoredMovies {
		finalMovies = append(
			finalMovies,
			scoredMovie.Movie,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": finalMovies,
	})
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

func (mc *MovieController) GetSuggestions(c *gin.Context) {

	search := c.Query("q")

	if len(strings.TrimSpace(search)) < 3 {
		c.JSON(http.StatusOK, []models.Movie{})
		return
	}

	search = utilities.NormalizeText(search)

	pipeline := mongo.Pipeline{
		//stage 1 = search with auto complete, typo tolerance, candidates ranking
		bson.D{{Key: "$search", Value: bson.D{
			{Key: "index", Value: "movies_search"},
			{Key: "autocomplete", Value: bson.D{
				{Key: "query", Value: search},
				{Key: "path", Value: "title"},
				{Key: "tokenOrder", Value: "sequential"},
				{Key: "fuzzy", Value: bson.D{
					{Key: "maxEdits", Value: 2},
				}},
			}},
		}}},
		//stage 2 = project, same as projection
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "imdb_id", Value: 1},
			{Key: "title", Value: 1},
			{Key: "poster_path", Value: 1},
		}}},
		//stage 3 = limit
		bson.D{{Key: "$limit", Value: 8}},
	}

	movieCollection := mc.db.Collection("movies")

	cursor, err := movieCollection.Aggregate(context.Background(), pipeline)

	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Error in fetching movies", gin.H{
			"endpoint": "/GetSuggestions",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch suggestions"})
		return
	}

	var movies []models.MovieSuggestion

	err = cursor.All(
		context.Background(),
		&movies,
	)

	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Failed to parse suggestions", gin.H{
			"endpoint": "/GetSuggestions",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch suggestions"})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (mc *MovieController) GetCachedHeroMovies() (map[string][]models.Movie, bool) {
	mc.HeroMovieCache.mu.Lock()
	defer mc.HeroMovieCache.mu.Unlock()

	if mc.HeroMovieCache.chachedData == nil {
		return nil, false
	}

	if time.Now().After(mc.HeroMovieCache.expiresAt) {
		return nil, false
	}

	return mc.HeroMovieCache.chachedData, true
}

func (mc *MovieController) SetCachedHeroMovies(heroMovies map[string][]models.Movie) {
	mc.HeroMovieCache.mu.Lock()
	defer mc.HeroMovieCache.mu.Unlock()

	mc.HeroMovieCache.chachedData = heroMovies
	mc.HeroMovieCache.expiresAt = time.Now().Add(10 * time.Minute)
}

func (mc *MovieController) GetHeroMovies(c *gin.Context) {
	//c.Header("Cache-Control", "public, max-age=300, s-maxage=600, stale-while-revalidate=60") //max-age = browser will keep this data for 300 secs, s-maxage = CDN will keep this data for 600 secs, stale-while-revalidate = keep the stale data alive while new data is refetched for 60 secs

	if cached, found := mc.GetCachedHeroMovies(); found {
		c.Header("X-Hero-Cache", "HIT")
		c.JSON(http.StatusOK, cached)
		return
	}

	c.Header("X-Hero-Cache", "MISS")

	movies := make(map[string][]models.Movie)

	var mu sync.Mutex

	var primaryWG sync.WaitGroup

	rankingQueries := []struct {
		key     string
		ranking string
	}{
		{"excellent", "Excellent"},
		{"good", "Good"},
		{"okay", "Okay"},
	}

	for _, query := range rankingQueries {

		primaryWG.Add(1)

		go func() {
			defer primaryWG.Done()

			result, err := mc.getMoviesByRanking(c, query.ranking, 10)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				movies[query.key] = []models.Movie{}
				return
			}

			movies[query.key] = result
		}()
	}

	primaryWG.Add(1)

	go func() {
		defer primaryWG.Done()

		result, err := mc.getNewMovies(c, 5, 3)

		mu.Lock()
		defer mu.Unlock()

		if err != nil {
			movies["newRelease"] = []models.Movie{}
			return
		}

		movies["newRelease"] = result
	}()

	primaryWG.Wait()

	err, genres := mc.fetchGenres(c)

	if err != nil {
		mc.dbLogger.Alerts("ERROR", "error fetching genres", gin.H{
			"endpoint": "/GetHeroMovies",
			"error":    err.Error(),
		})
		c.JSON(http.StatusOK, movies)
		return
	}

	var genreWG sync.WaitGroup

	for _, genre := range genres {

		genreWG.Add(1)

		go func() {
			defer genreWG.Done()

			result, err := mc.getMoviesByGenre(c, genre.GenreName, 10)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				movies[genre.GenreName] = []models.Movie{}
				return
			}

			movies[genre.GenreName] = result
		}()
	}

	genreWG.Wait()

	mc.SetCachedHeroMovies(movies)

	c.JSON(http.StatusOK, movies)
}

func (mc *MovieController) getMoviesByRanking(ctx context.Context, ranking string, limit int64) ([]models.Movie, error) {
	projection := heroMovieProjection()

	findOptions := options.Find().SetProjection(projection).SetSort(bson.M{"vote_average": -1})

	var movies []models.Movie

	moviesCollection := mc.db.Collection("movies")

	if limit != 0 {
		findOptions.SetLimit(limit)
	}

	cursor, err := moviesCollection.Find(ctx, bson.M{"ranking.ranking_name": ranking}, findOptions)
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "error fetching movies", gin.H{
			"endpoint": "/getMoviesByRanking",
			"error":    err.Error(),
		})
		return nil, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &movies); err != nil {
		mc.dbLogger.Alerts("ERROR", "error decoding movies", gin.H{
			"endpoint": "/getMoviesByRanking",
			"error":    err.Error(),
		})
		return nil, err
	}

	return movies, nil
}

func (mc *MovieController) getNewMovies(ctx context.Context, limit int64, yearOpts ...int) ([]models.Movie, error) {
	projection := heroMovieProjection()

	findOptions := options.Find().SetProjection(projection).SetSort(bson.M{"vote_average": -1})

	var movies []models.Movie

	moviesCollection := mc.db.Collection("movies")

	if limit != 0 {
		findOptions.SetLimit(limit)
	}

	year := 3
	if len(yearOpts) > 0 {
		year = yearOpts[0]
	}

	targetYear := time.Now().Year() - year
	minDateStr := fmt.Sprintf("%d-01-01", targetYear)

	filter := bson.M{
		"release_date": bson.M{"$gte": minDateStr},
	}

	cursor, err := moviesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "error fetching movies", gin.H{
			"endpoint": "/getNewMovies",
			"error":    err.Error(),
		})
		return nil, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &movies); err != nil {
		mc.dbLogger.Alerts("ERROR", "error decoding movies", gin.H{
			"endpoint": "/getNewMovies",
			"error":    err.Error(),
		})
		return nil, err
	}

	return movies, nil
}

func (mc *MovieController) getMoviesByGenre(ctx context.Context, genre string, limit int64) ([]models.Movie, error) {
	projection := heroMovieProjection()

	findOptions := options.Find().SetProjection(projection).SetSort(bson.M{"vote_average": -1})

	var movies []models.Movie

	moviesCollection := mc.db.Collection("movies")

	if limit != 0 {
		findOptions.SetLimit(limit)
	}

	cursor, err := moviesCollection.Find(ctx, bson.M{"genre.genre_name": genre}, findOptions)
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "error fetching movies", gin.H{
			"endpoint": "/getMoviesByGenre",
			"error":    err.Error(),
		})
		return nil, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &movies); err != nil {
		mc.dbLogger.Alerts("ERROR", "error decoding movies", gin.H{
			"endpoint": "/getMoviesByGenre",
			"error":    err.Error(),
		})
		return nil, err
	}

	return movies, nil
}

func (mc *MovieController) GetByCategory(c *gin.Context) {
	title := c.Param("title")

	var movies []models.Movie
	var err error

	categories := []string{"Excellent", "Good", "Okay"}
	nonGenreFound := slices.Contains(categories, title)

	if nonGenreFound {
		movies, err = mc.getMoviesByRanking(c, title, 0)
	} else if title == "newRelease" {
		movies, err = mc.getNewMovies(c, 0)
	} else {
		movies, err = mc.getMoviesByGenre(c, title, 0)
	}

	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Error in getting movies by category", gin.H{
			"endpoint": "/GetByCategory",
			"Category": title,
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (mc *MovieController) AddMovie(c *gin.Context) {

	role, err := utilities.GetRolefromContect(c)
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Error in getting role", gin.H{
			"endpoint": "/AddMovie",
			"error":    err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internal error"})
		return
	}

	if role != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be Admin to access this functionality"})
		return
	}

	var req struct {
		ImdbIDs string `json:"imdb_ids"`
	}

	if err := c.BindJSON(&req); err != nil {
		mc.dbLogger.Alerts("ERROR", "Bind error", gin.H{
			"endpoint": "/AddMovie",
			"error":    err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.ImdbIDs == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "imdb_ids required"})
		return
	}

	// Split input
	ids := strings.Split(req.ImdbIDs, ",")

	if len(ids) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too many IDs"})
		return
	}

	jobCollection := mc.db.Collection("jobs")

	added := []string{}
	skipped := []string{}

	for _, id := range ids {
		imdbID := strings.ToLower(strings.TrimSpace(id))

		if imdbID == "" {
			continue
		}

		// Basic validation
		if !utilities.IsValidImdbID(imdbID) {
			skipped = append(skipped, imdbID)
			continue
		}

		// Optional: prevent duplicate queueing
		// if queue.Exists(imdbID) {
		// 	skipped = append(skipped, imdbID)
		// 	continue
		// }

		job := models.Job{
			ImdbID:    imdbID,
			Status:    models.StatusPending,
			Attempts:  0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := jobCollection.InsertOne(context.Background(), job)
		if err != nil {
			continue
		}

		// queue.PushToQueue(imdbID)
		err = queue.EnqueueMovie(imdbID)

		if err != nil {
			mc.dbLogger.Alerts("ERROR", "Error in adding movies to queue", gin.H{
				"error":    err,
				"imdbID":   imdbID,
				"endpoint": "/AddMovie",
			})
			skipped = append(skipped, imdbID)
			continue
		}
		// queue.InspectMovieQueue(redisOpt)
		added = append(added, imdbID)
	}

	c.JSON(http.StatusOK, gin.H{
		"added":   added,
		"skipped": skipped,
	})
}

func (mc *MovieController) AddMovieViaAI(c *gin.Context) {
	role, err := utilities.GetRolefromContect(c)
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Error in getting role", gin.H{
			"endpoint": "/AddMovie",
			"error":    err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internal error"})
		return
	}

	if role != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be Admin to access this functionality"})
		return
	}

	var req struct {
		MovieDetails []struct {
			MovieName string `json:"movie_name"`
			MovieYear string `json:"movie_year"`
		} `json:"movie_details"`
	}

	if err := c.BindJSON(&req); err != nil {
		mc.dbLogger.Alerts("ERROR", "Bind error", gin.H{
			"endpoint": "/AddMovieViaAI",
			"error":    err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(req.MovieDetails) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movies required"})
		return
	}

	if len(req.MovieDetails) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too many IDs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movies added!!"})

	jobCollection := mc.db.Collection("jobs")

	skipped := []string{}
	added := []string{}

	for _, movies := range req.MovieDetails {

		movieName := strings.TrimSpace(movies.MovieName)
		movieYear := strings.TrimSpace(movies.MovieYear)
		if movieName == "" {
			continue
		}

		if movieYear != "" {
			if _, err := strconv.Atoi(movieYear); err != nil || len(movieYear) != 4 {
				movieYear = ""
			}
		}

		tmdbID, err := mc.GetTMDBIDByMovieName(movieName, movieYear, c)
		if err != nil || tmdbID == 0 {
			mc.dbLogger.Alerts("ERROR", "Error in fetching movie through movie name and year", gin.H{
				"error":    err,
				"tmdbID":   tmdbID,
				"endpoint": "/AddMovieViaAI",
			})
			skipped = append(skipped, movieName)
			continue
		}

		tmdbIDStr := strconv.Itoa(tmdbID)
		job := models.Job{
			ImdbID:    tmdbIDStr,
			Status:    models.StatusPending,
			Attempts:  0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     movieName,
		}

		_, err = jobCollection.InsertOne(context.Background(), job)
		if err != nil {
			continue
		}

		// queue.PushToQueue(imdbID)
		err = queue.EnqueueMovie(tmdbIDStr)

		if err != nil {
			mc.dbLogger.Alerts("ERROR", "Error in adding movies to queue", gin.H{
				"error":    err,
				"tmdbID":   tmdbID,
				"endpoint": "/AddMovieViaAI",
			})
			skipped = append(skipped, movies.MovieName)
			continue
		}
		// queue.InspectMovieQueue(redisOpt)
		added = append(added, movies.MovieName)
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"added":   added,
	// 	"skipped": skipped,
	// })
	
	mc.dbLogger.Log("ADD_MOVIE", "added and skipped movie", gin.H{
		"added":    added,
		"skipped":  skipped,
		"endpoint": "/AddMovieViaAI",
	})

}

func (mc *MovieController) GetTMDBIDByMovieName(movieName string, movieYear string, ctx context.Context) (int, error) {
	movieName = strings.TrimSpace(movieName)
	movieYear = strings.TrimSpace(movieYear)
	if movieName == "" {
		return 0, errors.New("movie name is required")
	}
	apiKey := strings.TrimSpace(os.Getenv("TMDB_SECRET_KEY"))
	if apiKey == "" {
		return 0, errors.New("TMDB_SECRET_KEY is not configured")
	}

	params := url.Values{}
	params.Set("query", movieName)
	params.Set("year", movieYear)
	params.Set("api_key", apiKey)
	apiURL := "https://api.themoviedb.org/3/search/movie?" + params.Encode()

	var result models.TMDBFindByMovieName
	err := utilities.WithRetry(5, ctx, func() error {
		return utilities.DoRequest(ctx, apiURL, &result, services.HttpClient)
	})
	if err != nil {
		return 0, err
	}

	if len(result.MovieResults) == 0 {
		mc.dbLogger.Alerts("WARN", "No Movie Results found", map[string]any{
			"endpoint":  "/GetTMDBIDByMovieName",
			"movieName": movieName,
			"movieYear": movieYear,
		})
		return 0, errors.New("no result found")
	}

	normalize := func(value string) string {
		return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), ""))
	}
	normalizedMovieName := normalize(movieName)
	for _, movie := range result.MovieResults {
		if movie.ID == 0 || normalize(movie.MovieName) != normalizedMovieName {
			continue
		}

		if movieYear == "" || strings.HasPrefix(movie.MovieYear, movieYear) {
			return movie.ID, nil
		}
	}

	return 0, errors.New("no matching movie found")
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

	var base_Prompt = os.Getenv("REVIEW_BASE_PROMPT")

	//Make this connection only once in main.go and pass it here to reduce time
	// client, err := genai.NewClient(ctx, option.WithAPIKey(GEMINI_API_KEY))

	client, err := GetGeminiClient(ctx)

	if err != nil {
		return "", 0, err
	}

	// model := client.GenerativeModel("gemini-2.5-flash") //gemini-3-flash-preview

	basePrompt := strings.Replace(base_Prompt, "{rankings}", allRankingNames, 1)

	fullPrompt := basePrompt + review

	// response, err := llm.Call(ctx, basePrompt+review)
	temp := float32(0)

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		genai.Text(fullPrompt),
		&genai.GenerateContentConfig{
			Temperature:      &temp,
			ResponseMIMEType: "text/plain",
		},
	)
	if err != nil {
		return "", 0, err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", 0, errors.New("empty response from Gemini")
	}

	response := strings.TrimSpace(resp.Text())

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
	var genreNames []models.Genre

	err, genreNames := mc.fetchGenres(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(http.StatusOK, genreNames)
}

func (mc *MovieController) fetchGenres(ctx context.Context) (error, []models.Genre) {
	genreCollection := mc.db.Collection("genres")

	var genreNames []models.Genre

	cursor, err := genreCollection.Find(ctx, bson.M{})
	if err != nil {
		mc.dbLogger.Alerts("ERROR", "Unable to fetch genres", gin.H{
			"endpoint": "/fetchGenres",
			"error":    err.Error(),
		})
		return err, nil
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var genreName models.Genre

		if err := cursor.Decode(&genreName); err != nil {
			mc.dbLogger.Alerts("WARN", "cursor decode error", gin.H{
				"endpoint": "/fetchGenres",
				"error":    err.Error(),
			})
		}

		genreNames = append(genreNames, genreName)
	}

	if err := cursor.Err(); err != nil {
		mc.dbLogger.Alerts("WARN", "cursor error", gin.H{
			"endpoint": "/fetchGenres",
			"error":    err.Error(),
		})
		return err, nil
	}

	return nil, genreNames
}

func GetGeminiClient(ctx context.Context) (*genai.Client, error) {
	var err error

	clientOnce.Do(func() {
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			err = errors.New("missing GEMINI_API_KEY")
			return
		}

		geminiClient, err = genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  apiKey,
			Backend: genai.BackendGeminiAPI,
		})
	})

	return geminiClient, err
}

func heroMovieProjection() bson.M {
	return bson.M{
		"title":         1,
		"imdb_id":       1,
		"poster_path":   1,
		"youtube_id":    1,
		"overview":      1,
		"release_date":  1,
		"runtime":       1,
		"is_adult":      1,
		"ranking":       1,
		"genre":         1,
		"tagline":       1,
		"backdrop_path": 1,
		"vote_average":  1,
	}
}
