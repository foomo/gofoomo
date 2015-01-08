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
	log.Println(response)
	return i
}

func (i *ImageInfo) getHeader(name string) []string {
	h, ok := i.Header[name]
	if !ok {
		h, _ = i.Header[strings.ToLower(name)]
	}
	return h
}

type Cache struct {
	Directory          map[string]*ImageInfo
	Foomo              *foomo.Foomo
	foomoSessionCookie *http.Cookie
	client             *http.Client
}

func NewCache(f *foomo.Foomo) *Cache {
	c := new(Cache)
	c.Foomo = f
	c.Directory = make(map[string]*ImageInfo)
	c.client = http.DefaultClient
	return c
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
		info = c.getImage(request, cookie)
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
	sessionCookie := getCookieByName(res.Cookies(), "foomoSessionTest")
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
		log.Println("requesting ", request.URL.String(), foomoMediaClientInfoCookie.String())
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
					log.Println("coul not parse expiration time", timeErr)
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
