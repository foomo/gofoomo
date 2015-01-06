package images

import (
	"io"
	"net/http"
	"os"
	"strings"
)

type Adaptive struct {
	Paths []string
}

type ImageCache interface {
	Get(incomingRequest *http.Request) *io.Reader
}

func NewAdaptive(paths []string) *Adaptive {
	a := new(Adaptive)
	a.Paths = paths
	return a
}

func (a *Adaptive) HandlesRequest(incomingRequest *http.Request) bool {
	for _, p := range a.Paths {
		if strings.HasPrefix(incomingRequest.URL.Path, p) {
			return true
		}
	}
	return false
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (a *Adaptive) ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request) {

}
