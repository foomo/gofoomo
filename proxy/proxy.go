package proxy

import (
	"github.com/foomo/gofoomo/foomo"
	"net/http"
	"net/http/httputil"
)

type Handler interface {
	HandlesRequest(incomingRequest *http.Request) bool
	ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request)
}

type Proxy struct {
	foomo        *foomo.Foomo
	reverseProxy *httputil.ReverseProxy
	handlers     []Handler
}

func NewProxy(f *foomo.Foomo) *Proxy {
	proxy := new(Proxy)
	proxy.foomo = f
	proxy.reverseProxy = httputil.NewSingleHostReverseProxy(proxy.foomo.URL)
	return proxy
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request) {
	for _, handler := range proxy.handlers {
		if handler.HandlesRequest(incomingRequest) {
			handler.ServeHTTP(w, incomingRequest)
			return
		}
	}
	incomingRequest.Host = proxy.foomo.URL.Host
	incomingRequest.URL.Opaque = incomingRequest.RequestURI
	proxy.reverseProxy.ServeHTTP(w, incomingRequest)
}

func (proxy *Proxy) AddHandler(handler Handler) {
	proxy.handlers = append(proxy.handlers, handler)
}
