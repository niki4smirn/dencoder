package server

import (
	"bytes"
	"dencoder/internal/storage"
	"dencoder/internal/tx"
	"dencoder/internal/utils"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

func (s *Server) Upload(w http.ResponseWriter, r *http.Request) error {
	// TODO: add logs
	logger := s.logger
	mpfile, h, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer mpfile.Close()

	all, err := io.ReadAll(mpfile)
	if err != nil {
		return err
	}
	logger.Infof("Client uploads file with size %v", len(all))

	link := fmt.Sprintf("upload/%s_%s%s", utils.FilenameWithoutExt(h.Filename), utils.RandSeq(10), filepath.Ext(h.Filename))

	transaction := tx.NewTx()

	transaction.Add(func(map[any]any) error {
		return storage.UploadVideo(s.cfg.S3BucketName, s.sess, link, bytes.NewReader(all), logger)
	}, func(map[any]any) error {
		return storage.DeleteVideo(s.cfg.S3BucketName, s.sess, link, logger)
	})

	transaction.Add(func(map[any]any) error {
		query := "INSERT INTO videos (filename, link) VALUES ($1, $2);"
		_, err = s.db.Exec(query, h.Filename, link)
		return err
	}, func(map[any]any) error {
		query := "DELETE FROM videos WHERE link = $1;"
		execRes, err := s.db.Exec(query, link)
		if err != nil {
			return err
		}

		rowsAffected, err := execRes.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected != 1 {
			logger.Errorf("Expected 1 row to be affected, actually %v", rowsAffected)
		}
		return nil
	})

	err, isFatal := transaction.Run()
	if err != nil {
		if isFatal {
			return fmt.Errorf("%w: %w", FatalErr, err)
		} else {
			return err
		}
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	err = tmpl.ExecuteTemplate(w, "video-list-element", VideoInfo{Filename: h.Filename, Link: link})
	if err != nil {
		return err
	}
	return nil
}
