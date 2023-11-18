package handlers

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
)

func readAll(reader io.ReadCloser) ([]byte, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("invalid reader: %s", err)
	}
	if err := reader.Close(); err != nil {
		return nil, fmt.Errorf("coudn't close reader: %s", err)
	}
	return body, nil
}

func readGzipBody(reader io.ReadCloser) ([]byte, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("invalid gzip reader: %s", err)
	}
	body, err := readAll(gzipReader)
	if err != nil {
		return nil, err
	}
	if err := reader.Close(); err != nil {
		return nil, fmt.Errorf("coudn't close gzip reader: %s", err)
	}
	return body, nil
}

// ReadBody - метод чтения тела API запроса.
// поддержка сжатия данных gzip.
func ReadBody(r *http.Request) ([]byte, error) {
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		log.Println("Gzip content")
		return readGzipBody(r.Body)
	default:
		return readAll(r.Body)
	}
}
