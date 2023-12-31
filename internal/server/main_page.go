package server

import (
	"html/template"
	"net/http"
)

type VideoInfo struct {
	id       int
	Link     string
	Filename string
}

type VideosList struct {
	Videos []VideoInfo
}

func (s *Server) MainPage(w http.ResponseWriter, r *http.Request) error {
	logger := s.logger
	// TODO: add context
	logger.Infof("Showing main page")
	rows, err := s.db.Query("SELECT * FROM videos;")
	if err != nil {
		return err
	}
	defer rows.Close()

	data := VideosList{}
	for rows.Next() {
		var info VideoInfo
		if err = rows.Scan(&info.id, &info.Filename, &info.Link); err != nil {
			return err
		}
		data.Videos = append(data.Videos, info)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}
