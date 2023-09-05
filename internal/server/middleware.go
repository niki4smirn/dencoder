package server

import "net/http"

func WithLogAndErr(handler func(http.ResponseWriter, *http.Request, *Logger) error, logger *Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r, logger); err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
