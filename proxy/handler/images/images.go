package images

import (
	"errors"
	"github.com/foomo/gofoomo/foomo"
	"net/http"
	"os"
	"strings"
	"time"
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

	info, err := a.Cache.Get(incomingRequest, a.BreakPoints)
	if err != nil {
		panic(err)
	}
	if info == nil {
		panic(errors.New("could not get image"))
	} else {
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
		w.Header().Set("Expires", time.Now().Add(time.Hour*24*30).Format(http.TimeFormat))
		//writeHeaders(w, info)
		http.ServeContent(w, incomingRequest, file.Name(), fileInfo.ModTime(), file)
		/*
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
			}*/
	}
}

func writeHeaders(w http.ResponseWriter, info *ImageInfo) {
	for key, values := range info.Header {
		if key == "Set-Cookie" {
			continue
		}

		/*if key == "Expires" {
			//force browser cache for 1 year
			w.Header().Add(key, time.Now().Add(time.Hour*24*365).Format(http.TimeFormat))

		}*/
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

}
