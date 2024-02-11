package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
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

func WithTimeout(timeout time.Duration, next func(w http.ResponseWriter, req *http.Request) error) func(w http.ResponseWriter, req *http.Request) error {
	return func(w http.ResponseWriter, req *http.Request) error {
		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		req = req.WithContext(ctx)

		err := next(w, req)
		if err != nil {
			return err
		}

		if err = ctx.Err(); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
                // TODO: handle properly
				w.Write([]byte("Timed out"))
				return err
			}
		}
		return nil
	}
}
