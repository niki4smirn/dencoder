package server

import (
	"dencoder/internal/storage"
	"fmt"
	"net/http"
)

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) error {
	logger := s.logger
	filename := r.URL.Query().Get("link")
	logger.Debugf("Removing %s", filename)
	if filename == "" {
		return fmt.Errorf("filename is not provided")
	}

	query := "DELETE FROM videos WHERE link = $1;"
	execRes, err := s.db.Exec(query, filename)
	if err != nil {
		return err
	}

	rowsAffected, err := execRes.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	if rowsAffected > 1 {
		logger.Errorf("Expected 1 row to be affected, actually %v", rowsAffected)
	}

	err = storage.DeleteVideo(s.cfg.S3BucketName, filename, logger)
	if err != nil {
		return err
	}

	return nil
}
