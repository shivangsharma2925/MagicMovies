package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	dblogger "github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/logger"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/models"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WatchListController struct {
	db       *database.MongoDB
	dbLogger *dblogger.DBLogger
}

func NewWatchListController(db *database.MongoDB, dbLogger *dblogger.DBLogger) *WatchListController {
	return &WatchListController{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (wc *WatchListController) NewWatchList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	userId, err := utilities.GetUserIdfromContext(c)
	if err != nil {
		wc.dbLogger.Alerts("ERROR", "Unable to get user id from context", gin.H{
			"endpoint": "/NewWatchList",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again!"})
		return
	}

	watchlistCollection := wc.db.Collection("watchlist")

	count, err := watchlistCollection.CountDocuments(ctx, bson.M{"user_id": userId})
	if err != nil {
		wc.dbLogger.Alerts("ERROR", "Error in executing db count query", gin.H{
			"endpoint": "/NewWatchList",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again!"})
		return
	}

	if count >= 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum limit of 5 watchlists reached"})
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		wc.dbLogger.Alerts("ERROR", "Invalid request for creating watch list", gin.H{
			"endpoint": "/NewWatchList",
			"error":    err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Watchlist name is required"})
		return
	}

	words := strings.Fields(name)
	var normalizedName string
	for _, word := range words {
		if len(word) == 0 {
			continue
		}
		normalizedName += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:]) + " "
	}
	normalizedName = strings.TrimSpace(normalizedName)

	watchList := models.WatchList{
		UserId:    userId,
		Name:      normalizedName,
		ColorIdx:  int(count) % 5,
		Movies:    []models.MoviesList{},
		CreatedAt: time.Now().UTC(),
	}

	_, err = watchlistCollection.InsertOne(ctx, watchList)
	if err != nil {
		wc.dbLogger.Alerts("ERROR", "Error in executing db insert query", gin.H{
			"endpoint": "/NewWatchList",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "New WatchList Added"})

}

func (wc *WatchListController) GetWatchLists(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	userId, err := utilities.GetUserIdfromContext(c)
	if err != nil {
		wc.dbLogger.Alerts("ERROR", "Unable to get user id from context", gin.H{
			"endpoint": "/GetWatchLists",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again!"})
		return
	}

	watchlistCollection := wc.db.Collection("watchlist")

	cursor, err := watchlistCollection.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		wc.dbLogger.Alerts("ERROR", "Error in executing db select query", gin.H{
			"endpoint": "/GetWatchLists",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again!"})
		return
	}

	var watchLists []models.WatchList

	err = cursor.All(ctx, &watchLists)
	if err != nil {
		wc.dbLogger.Alerts("ERROR", "Error in decoding watch lists", gin.H{
			"endpoint": "/GetWatchLists",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again!"})
		return
	}

	if watchLists == nil {
		watchLists = []models.WatchList{}
	}

	c.JSON(http.StatusOK, watchLists)

}

func (wc *WatchListController) RenameWatchLists(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	watchListIdStr := c.Param("watchlist_id")

	watchListId, err := primitive.ObjectIDFromHex(watchListIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Watchlist ID format"})
		return
	}

	var UserRequest struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&UserRequest); err != nil {
		wc.dbLogger.Alerts("ERROR", "Invalid request for updating name of watch list", gin.H{
			"endpoint": "/RenameWatchLists",
			"error":    err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	updatedName := strings.TrimSpace(UserRequest.Name)
	if len(updatedName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty name is not valid"})
		return
	}

	watchlistCollection := wc.db.Collection("watchlist")

	filter := bson.M{"_id": watchListId}
	update := bson.M{"$set": bson.M{"name": updatedName}}

	result, err := watchlistCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		wc.dbLogger.Alerts("ERROR", "Error in executing db update query", gin.H{
			"endpoint":     "/RenameWatchLists",
			"new_name":     updatedName,
			"watchlist_id": watchListIdStr,
			"error":        err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again!"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "WatchList does not exist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Name updated successfully"})

}

func (wc *WatchListController) DeleteWatchList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	watchListIdStr := c.Param("watchlist_id")

	watchListId, err := primitive.ObjectIDFromHex(watchListIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Watchlist ID format"})
		return
	}

	watchlistCollection := wc.db.Collection("watchlist")

	filter := bson.M{"_id": watchListId}

	result, err := watchlistCollection.DeleteOne(ctx, filter)
	if err != nil {
		wc.dbLogger.Alerts("ERROR", "Error in executing db delete query", gin.H{
			"endpoint":     "/DeleteWatchLists",
			"watchlist_id": watchListIdStr,
			"error":        err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again!"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "WatchList does not exist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "WatchList Deleted successfully"})

}

func (wc *WatchListController) AddMovietoWatchList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	watchListIdStr := c.Param("watchlist_id")

	watchListId, err := primitive.ObjectIDFromHex(watchListIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Watchlist ID format"})
		return
	}

	userID, err := utilities.GetUserIdfromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var req models.MoviesList

	if err := c.ShouldBindJSON(&req); err != nil {
		wc.dbLogger.Alerts("ERROR", "Invalid request for adding movie to watch list", gin.H{
			"endpoint": "/AddMovietoWatchList",
			"error":    err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	req.AddedAt = time.Now()
	req.UpdatedAt = time.Now()

	watchlistCollection := wc.db.Collection("watchlist")

	filter := bson.M{
		"_id":            watchListId,
		"user_id":        userID,
		"movies.imdb_id": bson.M{"$ne": req.ImdbID}, //ne = not equal
	}

	// filter = bson.M{
	// 	"_id": watchListId,
	// 	"movies.imdb_id": bson.M{
	// 		"nin": []string{req.ImdbID}, // nin = not in, nin expects the string to compare with
	// 	},
	// }

	// $addToSet to prevent duplicates automatically only if there are only static fields in the collection, any dynamic field will fail
	// update := bson.M{"$addToSet": bson.M{"movies": req}}

	update := bson.M{"$push": bson.M{"movies": req}} //Used "push" because movies is an array in watchlist collection

	result, err := watchlistCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		wc.dbLogger.Alerts("ERROR", "Error in executing db update query", gin.H{
			"endpoint":     "/AddMovietoWatchList",
			"watchlist_id": watchListIdStr,
			"error":        err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again!"})
		return
	}

	// If no document was modified, it means either the watchlist doesn't exist OR it's a duplicate
	if result.ModifiedCount == 0 {
		var existing models.WatchList
		errCount := watchlistCollection.FindOne(ctx, bson.M{"_id": watchListId}).Decode(&existing)

		if errCount == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "WatchList does not exist"})
			return
		}

		// If the watchlist exists but nothing updated, it was filtered out because of a duplicate ID
		c.JSON(http.StatusConflict, gin.H{"error": "Movie is already in this watchlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie added successfully"})

}

func (wc *WatchListController) DeleteMoviefromWatchList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	watchListIdStr := c.Param("watchlist_id")
	movieId := c.Param("movie_id")

	watchListId, err := primitive.ObjectIDFromHex(watchListIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Watchlist ID format"})
		return
	}

	userID, err := utilities.GetUserIdfromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	watchlistCollection := wc.db.Collection("watchlist")

	filter := bson.M{"_id": watchListId, "user_id": userID}
	update := bson.M{
		"$pull": bson.M{
			"movies": bson.M{"imdb_id": movieId},
		},
	}

	result, err := watchlistCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		wc.dbLogger.Alerts("ERROR", "Error in executing db pull query", gin.H{
			"endpoint":     "/DeleteMoviefromWatchList",
			"watchlist_id": watchListIdStr,
			"movie_id":     movieId,
			"error":        err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again!"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "WatchList does not exist"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie is not present in this watchlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie removed successfully"})

}

func (wc *WatchListController) ToggleWatch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var request struct {
		Watched *bool `json:"watched" binding:"required"`
	}

	watchListIdStr := c.Param("watchlist_id")
	movieId := c.Param("movie_id")

	watchListId, err := primitive.ObjectIDFromHex(watchListIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Watchlist ID format"})
		return
	}

	userID, err := utilities.GetUserIdfromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	err = c.ShouldBindJSON(&request)
	if err != nil {
		wc.dbLogger.Alerts("ERROR", "Error in binding the request", gin.H{
			"endpoint":     "/ToggleWatch",
			"watchlist_id": watchListIdStr,
			"movie_id":     movieId,
			"error":        err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "watched value is required"})
		return
	}

	watchlistCollection := wc.db.Collection("watchlist")

	filter := bson.M{"_id": watchListId, "user_id": userID, "movies.imdb_id": movieId}

	update := bson.M{
		"$set": bson.M{
			"movies.$.watched":    *request.Watched,
			"movies.$.updated_at": time.Now(),
		},
	}

	result, err := watchlistCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		wc.dbLogger.Alerts("ERROR", "Error in executing db update query", gin.H{
			"endpoint":     "/ToggleWatch",
			"watchlist_id": watchListIdStr,
			"movie_id":     movieId,
			"error":        err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, try again!"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Watchlist or movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie watched status updated successfully",
		"watched": request.Watched,
	})

}
