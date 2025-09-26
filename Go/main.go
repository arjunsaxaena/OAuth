package main

import (
	"log"

	"oauth/Go/common/config"
	"oauth/Go/common/connectors"
	"oauth/Go/common/logger"
)

func main() {
    // Initialize logger for this service
    logger.InitLogger(logger.OAUTH)

    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    // Create Postgres connection pool
    db := connectors.CreatePostgresSession(cfg.Database.DBURL)
    defer db.Close()

    select {}
}


