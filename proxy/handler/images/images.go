package images

import (
	"github.com/foomo/gofoomo/foomo"
	"io"
	"log"
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

func NewAdaptive(paths []string, breakPoints []int64, f *foomo.Foomo) *Adaptive {
	a := new(Adaptive)
	a.Paths = paths
	a.Cache = NewCache(f)
	a.BreakPoints = breakPoints
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
	file := a.Cache.Get(incomingRequest, a.BreakPoints)
	i, err := os.Stat(file.Name())
	log.Println("i, err", i, err)
	if file == nil {
		// dummy image
		w.Write([]byte("WTF"))
	} else {
		io.Copy(w, file)
		defer file.Close()
	}
}

func getExpiresFormattedTime() string {
	return time.Now().AddDate(0, 0, 2).Format("Mon, 02 Jan 2006 15:04:05 MST")
}
