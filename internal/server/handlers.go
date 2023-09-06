package server

import (
	"bytes"
	"dencoder/internal/utils"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

const batchSize = 1024 * 1024

type httpRange struct {
	left  uint64
	right uint64
}

func (r httpRange) len() uint64 {
	return r.right - r.left + 1
}

func parseRange(rangeStr string, fileSize uint64) (httpRange, error) {
	split := strings.Split(rangeStr, "-")
	if len(split) != 2 {
		return httpRange{}, fmt.Errorf("unexpected range header part: %v", rangeStr)
	}

	// suffix case
	if split[0] == "" {
		right, err := strconv.ParseUint(split[1], 10, 64)
		if err != nil {
			return httpRange{}, err
		}
		if fileSize < right {
			return httpRange{}, fmt.Errorf("suffix len (%v) > file size (%v)", right, fileSize)
		}
		left := fileSize - right
		return httpRange{left, right}, err
	}

	left, err := strconv.ParseUint(split[0], 10, 64)
	if err != nil {
		return httpRange{}, err
	}

	var right uint64
	if split[1] == "" {
		right = min(fileSize, left+batchSize) - 1
	} else {
		var err error
		right, err = strconv.ParseUint(split[1], 10, 64)
		if err != nil {
			return httpRange{}, err
		}
	}

	return httpRange{left, right}, nil
}

func parseRanges(rangesStr string, fileSize uint64) ([]httpRange, error) {
	prefix := "bytes="
	rangesStr, found := strings.CutPrefix(rangesStr, prefix)
	if !found {
		return []httpRange{}, fmt.Errorf("unexpected range header: %v", rangesStr)
	}

	split := strings.Split(rangesStr, ", ")

	res := make([]httpRange, 0)
	for _, curRangeStr := range split {
		parsedRange, err := parseRange(curRangeStr, fileSize)
		if err != nil {
			return nil, err
		}
		res = append(res, parsedRange)
	}
	return res, nil
}

func ServeVideo(content []byte, w http.ResponseWriter, r *http.Request, logger *Logger) error {
	// TODO: add logs
	logger.Infof("Serving video")
	fReader := bytes.NewReader(content)
	fSize := uint64(len(content))

	rangeHeader := r.Header.Get("Range")
	var contentRange httpRange

	if rangeHeader == "" {
		contentRange = httpRange{0, min(fSize, batchSize) - 1}
	} else {
		contentRanges, err := parseRanges(rangeHeader, fSize)
		if err != nil {
			return err
		}
		// WARNING: probably something criminal here
		contentRange = contentRanges[0]
	}

	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", strconv.FormatUint(contentRange.len(), 10))
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %v-%v/%v", contentRange.left, contentRange.right, fSize))
	w.Header().Set("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusPartialContent)

	_, err := fReader.Seek(int64(contentRange.left), io.SeekStart)
	if err != nil {
		return err
	}
	written, err := io.CopyN(w, fReader, int64(contentRange.len()))
	if written != int64(contentRange.len()) && err != io.EOF {
		return err
	}

	return nil
}

func Download(w http.ResponseWriter, r *http.Request, logger *Logger) error {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		return fmt.Errorf("filename is not provided")
	}

	// Use cache
	content, err := DownloadVideo(GetS3Bucket(r.Context()).Name, filename, logger)
	if err != nil {
		return err
	}

	return ServeVideo(content, w, r, logger)
}

func MainPage(w http.ResponseWriter, r *http.Request, logger *Logger) error {
	// TODO: add logs
	tmpl := template.Must(template.ParseFiles("index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		return err
	}
	return nil
}

func Upload(w http.ResponseWriter, r *http.Request, logger *Logger) error {
	// TODO: add logs
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

	filename := fmt.Sprintf("upload/%s_%s%s", utils.FilenameWithoutExt(h.Filename), utils.RandSeq(10), filepath.Ext(h.Filename))
	err = UploadVideo(GetS3Bucket(r.Context()).Name, filename, bytes.NewReader(all), logger)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	_, err = w.Write([]byte(filename))
	if err != nil {
		return err
	}

	return nil
}
