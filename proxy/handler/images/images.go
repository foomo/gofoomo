package images

import (
	"errors"

	"net/http"
	"os"
	"strings"
	"time"

	"github.com/foomo/gofoomo/foomo"
)

// Adaptive adaptive images
type Adaptive struct {
	Paths       []string
	BreakPoints []int64
	Cache       *Cache
}

// NewAdaptive constructor
func NewAdaptive(paths []string, f *foomo.Foomo) *Adaptive {
	a := new(Adaptive)
	a.Paths = paths
	a.Cache = NewCache(f)
	a.BreakPoints = getBreakPoints(f)
	return a
}

// HandlesRequest request handler interface implementation
func (a *Adaptive) HandlesRequest(incomingRequest *http.Request) bool {
	for _, p := range a.Paths {
		if strings.HasPrefix(incomingRequest.URL.Path, p) {
			return true
		}
	}
	return false
}

func (a *Adaptive) ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request) {
	info, err := a.Cache.Get(incomingRequest, a.BreakPoints)
	if err != nil {
		panic(err)
	}
	if info == nil {
		panic(errors.New("could not get image"))
	}
	// 304 handling
	file, err := os.Open(info.Filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()

	if err != nil {
		panic(err)
	}

	w.Header().Set("Expires", time.Now().Add(time.Hour*24*7).Format(http.TimeFormat))
	http.ServeContent(w, incomingRequest, file.Name(), fileInfo.ModTime(), file)
}

func writeHeaders(w http.ResponseWriter, info *ImageInfo) {
	for key, values := range info.Header {
		if key == "Set-Cookie" {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
}
