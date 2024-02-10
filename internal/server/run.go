package server

import (
	"database/sql"
	"dencoder/internal/config"
	"dencoder/internal/logging"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

type Logger = logging.Logger

func Run(cfg *config.Config, logger *Logger, db *sql.DB, sess *session.Session) error {
	// TODO: health check before run server (i.e. pgx and s3 consistency)
	router := chi.NewRouter()
	cache := &VideosCache{cache_data: expirable.NewLRU[string, VideoProvider](
		cfg.VideoCache.Size, nil, cfg.VideoCache.TTL,
	)}
	srv := Server{db, &cfg.ServerConfig, logger, sess, cache}
	// TODO: use context middleware (don't forget to use ctx in handler)
	router.Get("/get", WithErr(srv.ShowVideo, logger))
	router.Get("/delete", WithErr(srv.Delete, logger))
	router.Get("/", WithErr(srv.MainPage, logger))
	router.Post("/", WithErr(srv.Upload, logger))

	logger.Infof("Listening http://localhost:%v", cfg.ServerConfig.Port)
	return http.ListenAndServe(fmt.Sprintf(":%v", cfg.ServerConfig.Port), router)
}

type Server struct {
	db     *sql.DB
	cfg    *config.ServerConfig
	logger *Logger
	sess   *session.Session
	vCache *VideosCache
}
