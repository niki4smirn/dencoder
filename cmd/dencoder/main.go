package main

import (
	"database/sql"
	"dencoder/internal/config"
	"dencoder/internal/logging"
	"dencoder/internal/server"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("cannot load config: %w", err))
	}

	logger, err := logging.SetupLogger(cfg.Env)
	defer func(logger *logging.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)
	if err != nil {
		panic(fmt.Errorf("cannot setup logger: %w", err))
	}

	logger.Debug("config", cfg)

	pgxUser := os.Getenv("PGX_USER")
	pgxPass := os.Getenv("PGX_PASS")
	dbConnStr := fmt.Sprintf("postgresql://%s:%s@localhost/dencoder?sslmode=disable", pgxUser, pgxPass)

	db, err := sql.Open("pgx", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO: health check before run server (i.e. pgx and s3 consistency)
	if err := server.Run(cfg, logger, db); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
