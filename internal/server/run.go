package server

import (
	"dencoder/internal/config"
	"dencoder/internal/logging"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Logger = logging.Logger

func Run(cfg *config.Config, logger *Logger) error {
	router := chi.NewRouter()
	// use context middleware (don't forget to use ctx in handler)
	router.Get("/get", WithLogAndErr(Download, logger))
	router.Get("/", WithLogAndErr(MainPage, logger))
	router.Post("/", WithLogAndErr(Upload, logger))

	return http.ListenAndServe(fmt.Sprintf(":%v", cfg.HTTPConfig.Port), router)
}
