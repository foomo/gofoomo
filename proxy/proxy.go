package proxy

import (
	"github.com/foomo/gofoomo/foomo"
	"log"
	"net/http"
	"net/http/httputil"
)

type Handler interface {
	HandlesRequest(incomingRequest *http.Request) bool
	ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request)
}

type Listener interface {
	ListenServeHTTPStart(w http.ResponseWriter, incomingRequest *http.Request) http.ResponseWriter
	ListenServeHTTPDone(w http.ResponseWriter, incomingRequest *http.Request)
}

type Proxy struct {
	foomo        *foomo.Foomo
	ReverseProxy *httputil.ReverseProxy
	handlers     []Handler
	listeners    []Listener
	auth         *Auth
}

type ProxyServer struct {
	Proxy  *Proxy
	Config *Config
	Foomo  *foomo.Foomo
}

func NewProxy(f *foomo.Foomo) *Proxy {
	proxy := &Proxy{
		foomo: f,
	}
	proxy.ReverseProxy = httputil.NewSingleHostReverseProxy(proxy.foomo.URL)
	return proxy
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request) {
	if proxy.auth != nil && len(proxy.auth.Domain) > 0 && !proxy.foomo.BasicAuthForRequest(w, incomingRequest, proxy.auth.Domain, proxy.auth.Realm, "access denied") {
		return
	}
	for _, listener := range proxy.listeners {
		w = listener.ListenServeHTTPStart(w, incomingRequest)
	}
	for _, handler := range proxy.handlers {
		if handler.HandlesRequest(incomingRequest) {
			handler.ServeHTTP(w, incomingRequest)
			return
		}
	}
	incomingRequest.Host = proxy.foomo.URL.Host
	// incomingRequest.URL.Opaque = incomingRequest.RequestURI + incomingRequest.
	proxy.ReverseProxy.ServeHTTP(w, incomingRequest)

	for _, listener := range proxy.listeners {
		listener.ListenServeHTTPDone(w, incomingRequest)
	}

}

func (proxy *Proxy) AddHandler(handler Handler) {
	proxy.handlers = append(proxy.handlers, handler)
}

func (proxy *Proxy) AddListener(listener Listener) {
	proxy.listeners = append(proxy.listeners, listener)
}

func NewProxyServerWithConfig(filename string) (p *ProxyServer, err error) {
	config, err := ReadConfig(filename)
	if err != nil {
		return nil, err
	}
	return NewProxyServer(config)
}

func NewProxyServer(config *Config) (p *ProxyServer, err error) {
	proxyServer := new(ProxyServer)
	proxyServer.Config = config
	f, err := foomo.NewFoomo(config.Foomo.Dir, config.Foomo.RunMode, config.Foomo.Address)
	if err != nil {
		return nil, err
	}
	proxyServer.Foomo = f
	proxyServer.Proxy = NewProxy(proxyServer.Foomo)
	proxyServer.Proxy.auth = config.Server.Auth
	return proxyServer, nil
}

func (p *ProxyServer) ListenAndServe() error {
	c := p.Config.Server
	errorChan := make(chan error)
	if len(c.TLS.CertFile) > 0 && len(c.TLS.KeyFile) > 0 {
		log.Println("listening for https on", c.TLS.Address)
		go func() {
			errorChan <- http.ListenAndServeTLS(c.TLS.Address, c.TLS.CertFile, c.TLS.KeyFile, p.Proxy)
		}()
	}
	if len(c.Address) > 0 {
		log.Println("listening for http on", c.Address)
		go func() {
			errorChan <- http.ListenAndServe(c.Address, p.Proxy)
		}()
	}
	err := <-errorChan
	return err
}
