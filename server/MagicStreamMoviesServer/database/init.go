package database

import (
	"context"
	"log/slog"
	"time"
)

func InitDB(db *MongoDB, logger *slog.Logger) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call all index functions here
	if err := createTTLIndex(ctx, db); err != nil {
		logger.Error("TTL index creation failed", "error", err)
		return err
	}

	logger.Info("Database initialization completed")
	return nil
}
