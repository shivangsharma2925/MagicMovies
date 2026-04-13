package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	// "github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//this is not necessary now since main will now handle database connection calls
// var (
// 	clientInstance *mongo.Client
// 	clientOnce     sync.Once
// )

// func DBInstance() *mongo.Client {
// 	clientOnce.Do(func() {
// 		mongoURI := os.Getenv("MONGODB_URI")
// 		if mongoURI == "" {
// 			log.Fatal("MONGODB_URI not set")
// 		}
// 		...
// 	})
// 	return clientInstance
// }

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Initialize DB (called from main)
func NewMongoDB(ctx context.Context, uri, dbName string, logger *slog.Logger) (*MongoDB, error) {

	// it opens 50 simultaneous connection, keeps 10 connection ready all time, connection sitting idle for 5 mins are closed.
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(50).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(5 * time.Minute)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("mongo connect error: %w", err)
	}

	// Ping
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	logger.Info("Connected to MongoDB")

	return &MongoDB{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

// Clean way to get collection
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
