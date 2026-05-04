package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ImdbID      string             `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title       string             `bson:"title" json:"title" validate:"required,min=2,max=500"`
	PosterPath  string             `bson:"poster_path" json:"poster_path" validate:"required,url"`
	YouTubeID   string             `bson:"youtube_id" json:"youtube_id" validate:"required"`
	Genre       []Genre            `bson:"genre" json:"genre" validate:"required,dive"`
	AdminReview string             `bson:"admin_review" json:"admin_review" validate:"required"`
	Ranking     Ranking            `bson:"ranking" json:"ranking" validate:"required"`
}

type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required,min=2,max=100"`
}

type Ranking struct {
	RankingValue int    `bson:"ranking_value" json:"ranking_value" validate:"required"`
	RankingName  string `bson:"ranking_name" json:"ranking_name" validate:"required"` //oneof=Excellent Good Okay Bad Terible
}

type TMDBFindResponse struct {
	MovieResults []struct {
		ID int `json:"id"`
	} `json:"movie_results"`
}

type TMDBMovieDetails struct {
	Title      string  `json:"title"`
	PosterPath string  `json:"poster_path"`
	Rating     float64 `json:"vote_average"`
	Genres     []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
}

type TMDBVideoResponse struct {
	Results []struct {
		Key  string `json:"key"`
		Site string `json:"site"`
		Type string `json:"type"`
	} `json:"results"`
}

type MovieData struct {
	ImdbID    string
	Title     string
	Poster    string
	Genres    []Genre
	YoutubeID string
	Ranking   Ranking
}

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusProcessing JobStatus = "processing"
	StatusDone       JobStatus = "done"
	StatusFailed     JobStatus = "failed"
)

type Job struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ImdbID    string             `bson:"imdb_id" json:"imdb_id"`
	Status    JobStatus          `bson:"status" json:"status"`
	Attempts  int                `bson:"attempts" json:"attempts"`
	Error     string             `bson:"error,omitempty" json:"error,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
