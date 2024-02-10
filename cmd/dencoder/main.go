package main

import (
	"database/sql"
	"dencoder/internal/config"
	"dencoder/internal/logging"
	"dencoder/internal/server"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	pgxHost := os.Getenv("PGX_HOST")
	pgxPort := os.Getenv("PGX_PORT")
	pgxDatabase := os.Getenv("PGX_DATABASE")
	pgxUser := os.Getenv("PGX_USER")
	pgxPass := os.Getenv("PGX_PASS")
	dbConnStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		pgxUser, pgxPass, pgxHost, pgxPort, pgxDatabase)

	db, err := sql.Open("pgx", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Run(cfg, logger, db, sess); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
