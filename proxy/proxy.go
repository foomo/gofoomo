package proxy

import (
	"encoding/base64"
	"github.com/foomo/gofoomo/foomo"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Handler interface {
	HandlesRequest(incomingRequest *http.Request) bool
	ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request)
}

type Proxy struct {
	foomo        *foomo.Foomo
	reverseProxy *httputil.ReverseProxy
	handlers     []Handler
	auth         *Auth
}

type ProxyServer struct {
	Proxy  *Proxy
	Config *Config
	Foomo  *foomo.Foomo
}

func NewProxy(f *foomo.Foomo) *Proxy {
	proxy := new(Proxy)
	proxy.foomo = f
	proxy.reverseProxy = httputil.NewSingleHostReverseProxy(proxy.foomo.URL)
	return proxy
}

func (proxy *Proxy) forbidden(w http.ResponseWriter) {
	realm := strings.Replace(proxy.auth.Realm, "\"", "'", -1)
	w.Header().Set("Www-Authenticate", "Basic realm=\""+realm+"\"") //, encoding=\"UTF-8\"")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("access denied"))

}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request) {
	if proxy.auth != nil && len(proxy.auth.Domain) > 0 {
		authHeader := incomingRequest.Header.Get("Authorization")
		if len(authHeader) == 0 {
			proxy.forbidden(w)
			return
		}
		auth, base64DecodingErr := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		if base64DecodingErr != nil {
			proxy.forbidden(w)
			return
		}
		authParts := strings.Split(string(auth), ":")
		if len(authParts) != 2 {
			log.Println(string(auth), authParts, incomingRequest.Header)
			//panic(errors.New("malformed basic auth"))
			proxy.forbidden(w)
			return
		}
		if !proxy.foomo.BasicAuth(proxy.auth.Domain, authParts[0], authParts[1]) {
			proxy.forbidden(w)
			return
		}
	}
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
		go func() {
			errorChan <- http.ListenAndServeTLS(c.TLS.Address, c.TLS.CertFile, c.TLS.KeyFile, p.Proxy)

		}()
	}
	if len(c.Address) > 0 {
		go func() {
			errorChan <- http.ListenAndServe(c.Address, p.Proxy)
		}()
	}
	err := <-errorChan
	return err
}
