package handler

import (
	"net/http"

	"github.com/abbot/go-http-auth"
	"github.com/foomo/gofoomo/foomo"
	"github.com/foomo/gofoomo/proxy"
)

// AuthWrapper wraps proxy handlers with a basic authentication
type AuthWrapper struct {
	handler proxy.Handler
	foomo   *foomo.Foomo
	//	authDomain string
	denialHTML string
	//realm      string
	authenticator *auth.BasicAuth
}

// NewAuthWrapper wrap a proxy handler with basic auth
func NewAuthWrapper(f *foomo.Foomo, h proxy.Handler, authDomain string, realm string, denialHTML string) *AuthWrapper {
	secretProvider := auth.HtpasswdFileProvider(f.GetBasicAuthFilename(authDomain))
	return &AuthWrapper{
		handler: h,
		foomo:   f,
		//		realm:         realm,
		//		authDomain:    authDomain,
		denialHTML:    denialHTML,
		authenticator: auth.NewBasicAuthenticator(realm, secretProvider),
	}
}

// HandlesRequest am I responsible to handle that request
func (a *AuthWrapper) HandlesRequest(incomingRequest *http.Request) bool {
	// asking my underlying handler
	return a.handler.HandlesRequest(incomingRequest)
}

func (a *AuthWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user := a.authenticator.CheckAuth(r)
	if len(user) == 0 {
		a.authenticator.RequireAuth(w, r)
	}
	user = a.authenticator.CheckAuth(r)
	if len(user) == 0 {
		w.Write([]byte(a.denialHTML))
		return
	}

	a.handler.ServeHTTP(w, r)
}
