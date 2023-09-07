package server

import (
	"net/http"
)

func WithErr(handler func(http.ResponseWriter, *http.Request) error, logger *Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
