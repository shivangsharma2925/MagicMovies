package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WatchList struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId    string             `bson:"user_id" json:"user_id"`
	Name      string             `bson:"name" json:"name"`
	ColorIdx  int                `bson:"color_idx"  json:"color_idx"`
	Movies    []MoviesList       `bson:"movies" json:"movies"`
	CreatedAt time.Time          `bson:"time" json:"time"`
}

type MoviesList struct {
	ImdbID    string    `bson:"imdb_id" json:"imdb_id"`
	Title     string    `bson:"title" json:"title"`
	Poster    string    `bson:"poster"     json:"poster"`
	YouTubeID string    `bson:"youtube_id" json:"youtube_id"`
	Year      string    `bson:"year"       json:"year"`
	Watched   bool      `bson:"watched"    json:"watched"`
	AddedAt   time.Time `bson:"added_at" json:"added_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
