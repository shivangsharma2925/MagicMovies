package dblogger

import (
	"context"
	"time"

	"log/slog"

	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	"go.mongodb.org/mongo-driver/bson"
)

// example Alerts,
/* {
  "level": "ERROR",
  "message": "DB failure",
  "timestamp": "2026-04-01T10:00:00Z",
  "endpoint": "/login",
  "user_id": "123",
  "error": "connection timeout"
} */

// example Logs,
/* {
  "group": "account",
  "message": "Account created",
  "timestamp": "2026-04-01T10:00:00Z",
  "endpoint": "/login",
  "user_id": "123",
  "error": "connection timeout"
} */

type DBLogger struct {
	db     *database.MongoDB
	logger *slog.Logger
}

func NewDBLogger(db *database.MongoDB, logger *slog.Logger) *DBLogger {
	return &DBLogger{
		db:     db,
		logger: logger,
	}
}

func (l *DBLogger) Log(group, message string, meta map[string]any) {

	logsCollection := l.db.Collection("logs")

	logDoc := bson.M{
		"group":     group,
		"message":   message,
		"meta":      meta,
		"timestamp": time.Now(),
	}

	// Run async (important)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err := logsCollection.InsertOne(ctx, logDoc)
		if err != nil {
			l.logger.Error("Failed to write log to DB", "error", err)
		}
	}()
}

func (l *DBLogger) Alerts(level, message string, meta map[string]any) {

	alertsCollection := l.db.Collection("alerts")

	logDoc := bson.M{
		"level":     level,
		"message":   message,
		"meta":      meta,
		"timestamp": time.Now(),
	}

	// Run async (important)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err := alertsCollection.InsertOne(ctx, logDoc)
		if err != nil {
			l.logger.Error("Failed to write log to DB", "error", err)
		}
	}()
}