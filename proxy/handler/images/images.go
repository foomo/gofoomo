package images

import (
	"errors"
	"github.com/foomo/gofoomo/foomo"
	"io"
	//"log"
	"net/http"
	"os"
	"strings"
)

type Adaptive struct {
	Paths       []string
	BreakPoints []int64
	Cache       *Cache
}

func NewAdaptive(paths []string, f *foomo.Foomo) *Adaptive {
	a := new(Adaptive)
	a.Paths = paths
	a.Cache = NewCache(f)
	a.BreakPoints = getBreakPoints(f)
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

func (a *Adaptive) ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request) {
	info := a.Cache.Get(incomingRequest, a.BreakPoints)
	if info == nil {
		panic(errors.New("could not get image"))
	} else {
		// 304

		browserEtag := incomingRequest.Header.Get("If-None-Match")

		if browserEtag == info.Etag {
			w.WriteHeader(http.StatusNotModified)
			writeHeaders(w, info)
		} else {
			writeHeaders(w, info)
			file, err := os.Open(info.Filename)
			if err != nil {
				// dummy image ?!
				panic(errors.New("could not open image file " + info.Filename + " " + err.Error()))
			} else {
				io.Copy(w, file)
				defer file.Close()

			}
		}
	}
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
