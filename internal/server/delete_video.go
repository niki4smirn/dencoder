package server

import (
	"dencoder/internal/storage"
	"dencoder/internal/tx"
	"fmt"
	"net/http"
)

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) error {
	logger := s.logger
	link := r.URL.Query().Get("link")
	logger.Debugf("Removing %s", link)
	if link == "" {
		return fmt.Errorf("link is not provided")
	}

	transaction := tx.NewTx()

	transaction.Add(func(map[any]any) error {
		return storage.DeleteVideo(s.cfg.S3BucketName, link, logger)
	}, func(map[any]any) error {
		// it's strange to download video before deleting to just have a chance to recover it
		return fmt.Errorf("reverting video deletion is not supported")
	})

	const FilenameDataKey = "filename"
	transaction.Add(func(data map[any]any) error {
		query := "SELECT filename FROM videos WHERE link = $1"

		var filename string
		err := s.db.QueryRow(query, link).Scan(&filename)
		if err != nil {
			return err
		}

		data[FilenameDataKey] = filename

		query = "DELETE FROM videos WHERE link = $1;"
		execRes, err := s.db.Exec(query, link)
		if err != nil {
			return err
		}

		rowsAffected, err := execRes.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected > 1 {
			return fmt.Errorf("expected 1 row to be affected, actually %v", rowsAffected)
		}

		if rowsAffected == 0 {
			w.WriteHeader(http.StatusNotFound)
			return nil
		}

		return nil
	}, func(data map[any]any) error {
		filename, ok := data[FilenameDataKey].(string)
		if ok {
			return fmt.Errorf("data's filename field is not a string")
		}
		query := "INSERT INTO videos (filename, link) VALUES ($1, $2);"
		_, err := s.db.Exec(query, filename, link)
		return err
	})

	err, isFatal := transaction.Run()
	if err != nil {
		if isFatal {
			return fmt.Errorf("%w: %w", FatalErr, err)
		} else {
			return err
		}
	}

	return nil
}
