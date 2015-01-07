package images

import (
	"github.com/foomo/gofoomo/foomo"

	"log"
	"net/http"
	//"net/url"
	"os"
	//	"time"
)

type ImageInfo struct {
	Filename string
	Mimetype string
	Size     int64
}

type Cache struct {
	Files              map[string]string
	Foomo              *foomo.Foomo
	foomoSessionCookie *http.Cookie
	client             *http.Client
}

func NewCache(f *foomo.Foomo) *Cache {
	c := new(Cache)
	c.Foomo = f
	c.Files = make(map[string]string)
	c.client = http.DefaultClient
	return c
}

func (c *Cache) Get(request *http.Request, breakPoints []int64) *os.File {
	cookie := getFoomoMediaClientInfoCookie(request.Cookies(), breakPoints)
	key := cookie.String() + ":" + request.URL.Path
	filename, ok := c.Files[key]
	if ok == false {
		etag := c.getImage(request, cookie)
		if len(etag) > 0 {
			filename = c.Foomo.GetModuleCacheDir("Foomo.Media") + "/img-" + etag
			if fileExists(filename) {
				c.Files[key] = filename
			} else {
				return nil
			}
		}
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	return f
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (c *Cache) checkFoomoSessionCookie(res *http.Response) {
	sessionCookie := getCookieByName(res.Cookies(), FoomoMediaClientInfoCookieName)
	if sessionCookie != nil {
		if c.foomoSessionCookie == nil || (c.foomoSessionCookie != nil && c.foomoSessionCookie.Value != sessionCookie.Value) {
			log.Println("images.CheckFoomoSessionCookie: we have a session cookie", sessionCookie)
			c.foomoSessionCookie = sessionCookie
		}
	}
}

func (c *Cache) getImage(incomingRequest *http.Request, foomoMediaClientInfoCookie *http.Cookie) string {
	request, err := http.NewRequest("HEAD", incomingRequest.URL.String(), nil)
	//&url.URL{Host: incomingRequest.URL.Host, Scheme: incomingRequest.URL.Scheme, Opaque: incomingRequest.URL.Opaque}
	//requestTime := time.Now()
	if err != nil {
		return ""
	} else {
		request.AddCookie(foomoMediaClientInfoCookie)
		if c.foomoSessionCookie != nil {
			request.AddCookie(c.foomoSessionCookie)
		}
		// https://code.google.com/p/go/issues/detail?id=5684
		// http://godoc.org/net/url#URL ... Note that the Path field is stored in decoded form: /%47%6f%2f becomes /Go/. A
		request.URL.Opaque = incomingRequest.URL.Opaque
		request.URL.Host = c.Foomo.URL.Host
		request.URL.Scheme = c.Foomo.URL.Scheme
		imageServerResponse, err := c.client.Do(request)
		if imageServerResponse != nil {
			defer imageServerResponse.Body.Close()
		}
		if err != nil {
			if imageServerResponse != nil && imageServerResponse.StatusCode == http.StatusMovedPermanently {
				log.Fatalln("unexpected redirect")
				return ""
			} else {
				log.Fatalln("unexpected error", err)
				return ""
			}
		} else {
			c.checkFoomoSessionCookie(imageServerResponse)
			switch imageServerResponse.StatusCode {
			case http.StatusOK, http.StatusNotFound:
				if err != nil {
					log.Fatalln("unexpected error", err)
					return ""
				} else {
					return imageServerResponse.Header.Get("Etag")
				}
			default:
				log.Fatalln("unexpected reply", incomingRequest.URL, imageServerResponse, err)
				return ""
			}
		}
	}
}
