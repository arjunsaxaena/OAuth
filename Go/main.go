package main

import (
	"log"
	"net/http"

	"oauth/Go/api"
	"oauth/Go/common/config"
	"oauth/Go/common/connectors"
	"oauth/Go/common/logger"
	db "oauth/migrations/sqlc"
)

func main() {
    logger.InitLogger(logger.OAUTH)

    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    pgxPool := connectors.CreatePostgresSession(cfg.Database.DBURL)
    defer pgxPool.Close()

    store := db.New(pgxPool)
    handler := api.NewHandler(cfg, store)

    router := handler.Routes()

    logger.Info("Starting OAuth server on :4001")
    if err := http.ListenAndServe(":4001", router); err != nil {
        log.Fatalf("failed to start server: %v", err)
    }
}


