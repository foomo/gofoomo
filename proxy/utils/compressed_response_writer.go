package utils

import (
	"compress/gzip"
	//"io"
	"net/http"
	//"strings"
)

type CompressedResponseWriter struct {
	responseWriter http.ResponseWriter
	gz             *gzip.Writer
}

func NewCompressedResponseWriter(responseWriter http.ResponseWriter) *CompressedResponseWriter {
	crw := &CompressedResponseWriter{
		responseWriter: responseWriter,
		gz:             gzip.NewWriter(responseWriter),
	}
	return crw
}

// Header returns the header map that will be sent by WriteHeader.
// Changing the header after a call to WriteHeader (or Write) has
// no effect.
func (crw *CompressedResponseWriter) Header() http.Header {
	return crw.responseWriter.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (crw *CompressedResponseWriter) Write(bytes []byte) (int, error) {
	return crw.gz.Write(bytes)
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (crw *CompressedResponseWriter) WriteHeader(code int) {
	crw.responseWriter.WriteHeader(code)
}

func (crw *CompressedResponseWriter) Close() {
	crw.gz.Close()
}
