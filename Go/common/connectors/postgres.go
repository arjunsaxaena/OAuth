package connectors

import (
	"context"
	"oauth/Go/common/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Infinite iterator that returns the Postgres session
func CreatePostgresSession(dsn string) *pgxpool.Pool {
	count := 0
	for {
		db, err := pgxpool.New(context.Background(), dsn)
		if err == nil {
			// ACTUALLY TEST THE CONNECTION
			err = db.Ping(context.Background())
		}

		if err == nil {
			logger.Info("connected to postgres!")
			return db
		}

		count++
		logger.Error("unable to connect to postgres: %v", err)

		if count == 5 {
			logger.Info("retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			count = 0
		}
	}
}

