package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createTTLIndex(ctx context.Context, db *MongoDB) error {
	logsLollection := db.Collection("logs")
	AlertsCollection := db.Collection("alerts")

	logsIndex := mongo.IndexModel{
		Keys: bson.M{"timestamp": 1},
		Options: options.Index().
			SetExpireAfterSeconds(86400).
			SetName("logs_ttl_index"),
	}

	alertsIndex := mongo.IndexModel{
		Keys: bson.M{"timestamp": 1},
		Options: options.Index().
			SetExpireAfterSeconds(3600).
			SetName("alerts_ttl_index"),
	}

	_, err := AlertsCollection.Indexes().CreateOne(ctx, alertsIndex)
	_, err = logsLollection.Indexes().CreateOne(ctx, logsIndex)
	return err
}
