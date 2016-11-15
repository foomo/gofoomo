package proxy

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/foomo/gofoomo/foomo"
	"github.com/foomo/tlsconfig"
)

// Handler proxy handler
type Handler interface {
	HandlesRequest(incomingRequest *http.Request) bool
	ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request)
}

// Listener deprecated
type Listener interface {
	ListenServeHTTPStart(w http.ResponseWriter, incomingRequest *http.Request) http.ResponseWriter
	ListenServeHTTPDone(w http.ResponseWriter, incomingRequest *http.Request)
}

// Proxy foomo proxy
type Proxy struct {
	foomo         *foomo.Foomo
	ReverseProxy  *httputil.ReverseProxy
	handlers      []Handler
	listeners     []Listener
	auth          *Auth
	ServeHTTPFunc func(http.ResponseWriter, *http.Request)
}

// Server server for Proxy
type Server struct {
	Proxy     *Proxy
	Config    *Config
	TLSConfig *tls.Config
	Foomo     *foomo.Foomo
}

// NewProxy constructor
func NewProxy(f *foomo.Foomo) *Proxy {
	proxy := &Proxy{
		foomo: f,
	}
	proxy.ServeHTTPFunc = proxy.serveHTTP
	proxy.ReverseProxy = httputil.NewSingleHostReverseProxy(proxy.foomo.URL)
	return proxy
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request) {
	proxy.ServeHTTPFunc(w, incomingRequest)
}

func (proxy *Proxy) serveHTTP(w http.ResponseWriter, incomingRequest *http.Request) {
	if proxy.auth != nil && len(proxy.auth.Domain) > 0 {
		//&& !proxy.foomo.BasicAuthForRequest(w, incomingRequest, proxy.auth.Domain, proxy.auth.Realm, "access denied") {

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

// AddHandler add a handler
func (proxy *Proxy) AddHandler(handler Handler) {
	proxy.handlers = append(proxy.handlers, handler)
}

// AddListener deprecated
func (proxy *Proxy) AddListener(listener Listener) {
	proxy.listeners = append(proxy.listeners, listener)
}

// NewServerWithConfig constructor with a config file
func NewServerWithConfig(filename string) (p *Server, err error) {
	config, err := ReadConfig(filename)
	if err != nil {
		return nil, err
	}
	return NewServer(config)
}

// NewServer constructor with config struct
func NewServer(config *Config) (p *Server, err error) {
	p = &Server{
		Config: config,
	}
	f, err := foomo.NewFoomo(config.Foomo.Dir, config.Foomo.RunMode, config.Foomo.Address)
	if err != nil {
		return nil, err
	}
	p.Foomo = f
	p.Proxy = NewProxy(p.Foomo)
	p.Proxy.auth = config.Server.Auth
	p.TLSConfig = tlsconfig.NewServerTLSConfig(p.Config.Server.TLS.Mode)
	return p, nil
}

// ListenAndServe until things go bad, depending upon configuration this will\
// listen to http and https requests
func (p *Server) ListenAndServe() error {
	c := p.Config.Server
	errorChan := make(chan error)
	startedHTTPS := false
	startedHTTP := false
	if len(c.TLS.CertFile) > 0 && len(c.TLS.KeyFile) > 0 {
		log.Println("listening for https on", c.TLS.Address)
		go func() {
			tlsServer := &http.Server{
				Addr:      c.TLS.Address,
				Handler:   p.Proxy,
				TLSConfig: p.TLSConfig,
			}
			errorChan <- tlsServer.ListenAndServeTLS(c.TLS.CertFile, c.TLS.KeyFile)
		}()
		startedHTTPS = true
	}
	if len(c.Address) > 0 {
		log.Println("listening for http on", c.Address)
		go func() {
			errorChan <- http.ListenAndServe(c.Address, p.Proxy)
		}()
		startedHTTP = true
	}
	if !startedHTTP && !startedHTTPS {
		return errors.New("nothing to listen to")
	}
	err := <-errorChan
	return err
}
