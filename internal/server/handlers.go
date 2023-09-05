package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
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

func ServeVideo(filename string, w http.ResponseWriter, r *http.Request, logger *Logger) error {
	// TODO: add logs
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	fReader := bytes.NewReader(file)
	fSize := uint64(len(file))

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

	_, err = fReader.Seek(int64(contentRange.left), io.SeekStart)
	if err != nil {
		return err
	}
	written, err := io.CopyN(w, fReader, int64(contentRange.len()))
	if written != int64(contentRange.len()) && err != io.EOF {
		return err
	}

	return nil
}

func HandleGet(w http.ResponseWriter, r *http.Request, logger *Logger) error {
	return ServeVideo("/home/niki4smirn/Dev/go/dencoder/video.MP4", w, r, logger)
}
