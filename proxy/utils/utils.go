package utils

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func ServeCompressed(w http.ResponseWriter, incomingRequest *http.Request, writeCallback func(writer io.Writer) error) error {
	var writer io.Writer
	writer = w
	if strings.Contains(incomingRequest.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		writer = gz
	}
	return writeCallback(writer)

}

type Resource304 struct {
	Etag string
}

func Handle304(r *http.Request) {}
