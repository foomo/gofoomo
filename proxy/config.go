package proxy

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/foomo/tlsconfig"

	"gopkg.in/yaml.v2"
)

// Auth basic auth
type Auth struct {
	Domain string
	Realm  string
}

// TLS proxy tls config vo
type TLS struct {
	Mode     tlsconfig.TLSModeServer
	Address  string
	CertFile string
	KeyFile  string
}

// Config proxy configuration
type Config struct {
	// how should the proxy server run
	Server struct {
		Address string
		Auth    *Auth
		TLS     TLS
	}
	// where is foomo
	Foomo struct {
		// php server address like http://test.foomo
		Address string
		// test, development or production
		RunMode string
		// locally accessible directory for the server
		Dir string
	}
	// this is for you and your handlers
	AppOptions map[string]string
}

// ReadConfig from a file
func ReadConfig(filename string) (config *Config, err error) {
	config = &Config{}
	yamlBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlBytes, &config)
	if err != nil {
		return nil, err
	}
	if config.Server.TLS.Mode == "" {
		config.Server.TLS.Mode = tlsconfig.TLSModeServerDefault
	}
	switch config.Server.TLS.Mode {
	case tlsconfig.TLSModeServerDefault, tlsconfig.TLSModeServerLoose, tlsconfig.TLSModeServerStrict:
	default:
		return nil, errors.New("unknown server.tls.mode: " + string(config.Server.TLS.Mode) + " - must be one of: " + fmt.Sprintln([]tlsconfig.TLSModeServer{tlsconfig.TLSModeServerDefault, tlsconfig.TLSModeServerLoose, tlsconfig.TLSModeServerStrict}))
	}
	return config, nil
}
