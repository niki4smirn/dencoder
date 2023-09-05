package main

import (
	"dencoder/internal/config"
	"dencoder/internal/logging"
	"dencoder/internal/server"
	"fmt"
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

	if err := server.Run(cfg, logger); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
