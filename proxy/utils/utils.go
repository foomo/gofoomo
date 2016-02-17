package utils

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// ServeCompressed serve a request with a compressed response, if the client accepts it
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

// Resource304 vo
type Resource304 struct {
	Etag string
}

// Handle304 304 handler
func Handle304(r *http.Request) {}
