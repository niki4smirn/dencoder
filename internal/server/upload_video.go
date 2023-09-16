package server

import (
	"bytes"
	"dencoder/internal/storage"
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
	err = storage.UploadVideo(s.cfg.S3BucketName, link, bytes.NewReader(all), logger)
	if err != nil {
		return err
	}

	// maybe rollback s3 upload in case of db failure

	query := "INSERT INTO videos (filename, link) VALUES ($1, $2);"
	_, err = s.db.Exec(query, h.Filename, link)
	if err != nil {
		return err
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	err = tmpl.ExecuteTemplate(w, "video-list-element", VideoInfo{Filename: h.Filename, Link: link})
	if err != nil {
		return err
	}
	return nil
}
