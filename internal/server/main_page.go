package server

import (
	"html/template"
	"net/http"
)

func (s *Server) MainPage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	logger := s.logger
	logger.Infof("Showing main page")
	if s.vInfoCache.Videos == nil {

		rows, err := s.db.QueryContext(ctx, "SELECT * FROM videos;")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var info VideoInfo
			if err = rows.Scan(&info.id, &info.Filename, &info.Link); err != nil {
				return err
			}
			s.vInfoCache.Videos = append(s.vInfoCache.Videos, info)
		}
		if err = rows.Err(); err != nil {
			return err
		}
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	err := tmpl.Execute(w, s.vInfoCache)
	if err != nil {
		return err
	}
	return nil
}
