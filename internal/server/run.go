package server

import (
	"database/sql"
	"dencoder/internal/config"
	"dencoder/internal/logging"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Logger = logging.Logger

func Run(cfg *config.Config, logger *Logger, db *sql.DB) error {
	router := chi.NewRouter()
	srv := Server{db, &cfg.ServerConfig, logger}
	// use context middleware (don't forget to use ctx in handler)
	router.Get("/get", WithErr(srv.Download, logger))
	router.Get("/", WithErr(srv.MainPage, logger))
	router.Post("/", WithErr(srv.Upload, logger))

	return http.ListenAndServe(fmt.Sprintf(":%v", cfg.ServerConfig.Port), router)
}

type Server struct {
	db     *sql.DB
	cfg    *config.ServerConfig
	logger *Logger
}
