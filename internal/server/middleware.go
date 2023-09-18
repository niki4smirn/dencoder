package server

import (
	"errors"
	"fmt"
	"net/http"
)

var FatalErr = fmt.Errorf("fatal error")

func WithErr(handler func(http.ResponseWriter, *http.Request) error, logger *Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			if errors.Is(err, FatalErr) {
				logger.Fatal(err)
			} else {
				logger.Error(err)
			}
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
