package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ImdbID          string             `bson:"imdb_id" json:"imdb_id" validate:"required"`
	TmdbID          int                `bson:"tmdb_id" json:"tmdb_id" validate:"required"`
	Title           string             `bson:"title" json:"title" validate:"required,min=2,max=500"`
	TitleNormalized string             `bson:"title_normalized,omitempty" json:"-"`
	SearchTokens    []string           `bson:"search_tokens,omitempty" json:"-"`
	PosterPath      string             `bson:"poster_path" json:"poster_path" validate:"required,url"`
	BackdropPath    string             `bson:"backdrop_path" json:"backdrop_path"`
	YouTubeID       string             `bson:"youtube_id" json:"youtube_id" validate:"required"`
	Genre           []Genre            `bson:"genre" json:"genre" validate:"required,dive"`
	AdminReview     string             `bson:"admin_review" json:"admin_review" validate:"required"`
	Ranking         Ranking            `bson:"ranking" json:"ranking" validate:"required"`
	IsAdult         bool               `bson:"is_adult" json:"adult"`
	Overview        string             `bson:"overview" json:"overview"`
	ReleaseDate     string             `bson:"release_date" json:"release_date"`
	Tagline         string             `bson:"tagline" json:"tagline"`
	Runtime         int                `bson:"runtime" json:"runtime"`
	VoteAvg         float64            `bson:"vote_average" json:"vote_average"`
}

type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required,min=2,max=100"`
}

type Ranking struct {
	RankingValue int    `bson:"ranking_value" json:"ranking_value" validate:"required"`
	RankingName  string `bson:"ranking_name" json:"ranking_name" validate:"required"` //oneof=Excellent Good Okay Bad Terible
}

type MovieSuggestion struct {
	ImdbID     string `bson:"imdb_id" json:"imdb_id"`
	Title      string `bson:"title" json:"title"`
	PosterPath string `bson:"poster_path" json:"poster_path"`
}

// Manually giving search results - not using anywhere
type ScoredMovie struct {
	Movie Movie
	Score float64
}

type TMDBFindResponse struct {
	MovieResults []struct {
		ID int `json:"id"`
	} `json:"movie_results"`
}

type TMDBFindByMovieName struct {
	MovieResults []struct {
		ID        int    `json:"id"`
		MovieName string `json:"title"`
		MovieYear string `json:"release_date"`
	} `json:"results"`
}

type TMDBMovieDetails struct {
	ImdbID     string `json:"imdb_id"`
	Title      string `json:"title"`
	PosterPath string `json:"poster_path"`
	Genres     []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	IsAdult      bool    `json:"adult"`
	Overview     string  `json:"overview"`
	ReleaseDate  string  `json:"release_date"`
	Tagline      string  `json:"tagline"`
	Runtime      int     `json:"runtime"`
	BackdropPath string  `json:"backdrop_path"`
	VoteAvg      float64 `json:"vote_average"`
}

type TMDBVideoResponse struct {
	Results []struct {
		Key  string `json:"key"`
		Site string `json:"site"`
		Type string `json:"type"`
	} `json:"results"`
}

type MovieData struct {
	ImdbID       string
	TmdbID       int
	Title        string
	Poster       string
	Genres       []Genre
	YoutubeID    string
	Ranking      Ranking
	IsAdult      bool
	Overview     string
	ReleaseDate  string
	Tagline      string
	Runtime      int
	BackdropPath string
	VoteAvg      float64
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
	Title     string             `bson:"title" json:"title"`
	Status    JobStatus          `bson:"status" json:"status"`
	Attempts  int                `bson:"attempts" json:"attempts"`
	Error     string             `bson:"error,omitempty" json:"error,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
