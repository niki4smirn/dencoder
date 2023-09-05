package logging

import (
	"dencoder/internal/config"
	"fmt"
	"go.uber.org/zap"
)

type Logger = zap.SugaredLogger

func SetupLogger(env config.Env) (*Logger, error) {
	var logger *zap.Logger
	var sugaredLogger *Logger
	var err error
	switch env {
	case config.Dev:
		logger, err = zap.NewDevelopment()
	case config.Prod:
		logger, err = zap.NewProduction()
	default:
		return nil, fmt.Errorf("unexpected env value: %v", env)
	}
	if err != nil {
		return nil, err
	}
	sugaredLogger = logger.Sugar()
	return sugaredLogger, nil
}
