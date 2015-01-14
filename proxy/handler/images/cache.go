package images

import (
	"errors"
	"github.com/foomo/gofoomo/foomo"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ImageInfo struct {
	Filename   string
	Etag       string
	StatusCode int
	Header     map[string][]string
	Expires    int64
}

func NewImageInfo(response *http.Response) *ImageInfo {
	i := new(ImageInfo)
	if response != nil {
		i.StatusCode = response.StatusCode
		i.Header = response.Header
	}
	return i
}

func (i *ImageInfo) getHeader(name string) []string {
	h, ok := i.Header[name]
	if !ok {
		h, _ = i.Header[strings.ToLower(name)]
	}
	return h
}

type ImageRequest struct {
	Id                         string
	IncomingRequest            *http.Request
	FoomoMediaClientInfoCookie *http.Cookie
	DoneChannel                chan *ImageInfo
	ImageInfo                  *ImageInfo
}

func NewImageRequest(id string, incomingRequest *http.Request, foomoMediaClientInfoCookie *http.Cookie) *ImageRequest {
	r := new(ImageRequest)
	r.Id = id
	r.DoneChannel = make(chan *ImageInfo)
	r.IncomingRequest = incomingRequest
	r.FoomoMediaClientInfoCookie = foomoMediaClientInfoCookie
	return r
}

func (i *ImageRequest) execute(cache *Cache) {
	cache.RequestChannel <- i
	i.ImageInfo = <-i.DoneChannel

}

type Cache struct {
	Directory              map[string]*ImageInfo
	Foomo                  *foomo.Foomo
	foomoSessionCookie     *http.Cookie
	foomoSessionCookieName string
	client                 *http.Client
	RequestChannel         chan *ImageRequest

	//doneChannel        chan *ImageRequest
}

func NewCache(f *foomo.Foomo) *Cache {
	c := new(Cache)
	c.Foomo = f
	c.Directory = make(map[string]*ImageInfo)
	c.client = http.DefaultClient
	c.foomoSessionCookieName = getFoomoSessionCookieName(f)
	c.RequestChannel = make(chan *ImageRequest)
	//c.doneChannel = make(chan *ImageInfo)
	go c.runLoop()
	return c
}

func (c *Cache) runLoop() {
	pendingRequests := make(map[string][]*ImageRequest)
	doneChannel := make(chan *ImageRequest)
	for {
		select {
		case r := <-c.RequestChannel:
			// incoming request
			_, ok := pendingRequests[r.Id]
			if !ok {
				// that is a new one
				pendingRequests[r.Id] = []*ImageRequest{}
				go func() {
					r.ImageInfo = c.getImage(r.IncomingRequest, r.FoomoMediaClientInfoCookie)
					doneChannel <- r
				}()
			} else {
				log.Println("hang on")
			}
			pendingRequests[r.Id] = append(pendingRequests[r.Id], r)
		case done := <-doneChannel:
			requests, _ := pendingRequests[done.Id]
			for _, r := range requests {
				r.DoneChannel <- done.ImageInfo
			}
			delete(pendingRequests, done.Id)
		}
	}

}

func (c *Cache) Get(request *http.Request, breakPoints []int64) *ImageInfo {
	cookie := getFoomoMediaClientInfoCookie(request.Cookies(), breakPoints)
	key := cookie.String() + ":" + request.URL.Path
	info, ok := c.Directory[key]
	if ok && time.Now().Unix() > info.Expires {
		log.Println("that image expired - getting a new one", info.Expires, time.Now())
		ok = false
		info = nil
		delete(c.Directory, key)
	}
	if ok == false {
		imageRequest := NewImageRequest(key, request, cookie)
		imageRequest.execute(c)
		info = imageRequest.ImageInfo
		if len(info.Etag) > 0 {
			info.Filename = c.Foomo.GetModuleCacheDir("Foomo.Media") + "/img-" + info.Etag
			if fileExists(info.Filename) {
				c.Directory[key] = info
			} else {
				return nil
			}
		}
	}
	return info
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (c *Cache) checkFoomoSessionCookie(res *http.Response) {
	sessionCookie := getCookieByName(res.Cookies(), c.foomoSessionCookieName)
	if sessionCookie != nil {
		if c.foomoSessionCookie == nil || (c.foomoSessionCookie != nil && c.foomoSessionCookie.Value != sessionCookie.Value) {
			log.Println("images.CheckFoomoSessionCookie: we have a session cookie", sessionCookie)
			c.foomoSessionCookie = sessionCookie
		}
	}
}

func (c *Cache) getImage(incomingRequest *http.Request, foomoMediaClientInfoCookie *http.Cookie) *ImageInfo {
	request, err := http.NewRequest("HEAD", incomingRequest.URL.String(), nil)
	if err != nil {
		return NewImageInfo(nil)
	} else {
		request.AddCookie(foomoMediaClientInfoCookie)
		if c.foomoSessionCookie != nil {
			request.AddCookie(c.foomoSessionCookie)
		}
		request.URL.Opaque = incomingRequest.URL.Opaque
		request.URL.Host = c.Foomo.URL.Host
		request.URL.Scheme = c.Foomo.URL.Scheme
		imageServerResponse, err := c.client.Do(request)
		i := NewImageInfo(imageServerResponse)
		if imageServerResponse != nil {
			defer imageServerResponse.Body.Close()
		}
		if err != nil {
			if imageServerResponse != nil && imageServerResponse.StatusCode == http.StatusMovedPermanently {
				panic(errors.New("unexpected redirect"))
			} else {
				panic(errors.New("unexpected error " + err.Error()))
			}
		} else {
			i.StatusCode = imageServerResponse.StatusCode
			c.checkFoomoSessionCookie(imageServerResponse)
			switch i.StatusCode {
			case http.StatusOK, http.StatusNotFound:
				t, timeErr := time.Parse(time.RFC1123, imageServerResponse.Header.Get("Expires"))
				if timeErr == nil {
					i.Expires = t.Unix()
				} else {
					i.Expires = 0
					i.Expires = time.Now().Unix() + 3600
					log.Println("could not parse expiration time", timeErr)
				}
				if err != nil {
					panic(errors.New("unexpected error " + err.Error()))
				} else {
					i.Etag = imageServerResponse.Header.Get("Etag")
				}
			default:
				panic(errors.New("unexpected reply with status " + imageServerResponse.Status))
			}
		}
		return i
	}
}
