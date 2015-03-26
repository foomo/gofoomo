package handler

import (
	"github.com/foomo/gofoomo/foomo"
	"github.com/foomo/gofoomo/proxy"
	"log"
	"net/http"
)

type AuthWrapper struct {
	handler    proxy.Handler
	foomo      *foomo.Foomo
	authDomain string
	denialHTML string
	realm      string
}

func NewAuthWrapper(f *foomo.Foomo, h proxy.Handler, authDomain string, realm string, denialHTML string) *AuthWrapper {
	return &AuthWrapper{
		handler:    h,
		foomo:      f,
		realm:      realm,
		authDomain: authDomain,
		denialHTML: denialHTML,
	}
}

func (a *AuthWrapper) HandlesRequest(incomingRequest *http.Request) bool {
	return a.handler.HandlesRequest(incomingRequest)
}

func (a *AuthWrapper) ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request) {
	if !a.foomo.BasicAuthForRequest(w, incomingRequest, a.authDomain, "authenticate", "nope") {
		log.Println("auth wrapper access denied for", incomingRequest.RequestURI)
		return
	}
	a.handler.ServeHTTP(w, incomingRequest)
}
