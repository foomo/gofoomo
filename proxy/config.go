package proxy

import (
	"crypto/tls"
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

type Certificate struct {
	CertFile string
	KeyFile  string
}

// TLS proxy tls config vo
type TLS struct {
	Mode         tlsconfig.TLSModeServer
	ForceTLS     bool
	Address      string
	CertFile     string
	KeyFile      string
	Certificates []Certificate
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

func (c *Config) getCertificates() (certificates []tls.Certificate, err error) {
	if !c.hasCertificates() {
		err = errors.New("i do not have any certificates")
		return
	}
	if (c.Server.TLS.CertFile != "" || c.Server.TLS.KeyFile != "") && len(c.Server.TLS.Certificates) > 0 {
		err = errors.New("you can not mix .Certificates and .CertFile and .KeyFile - choose one")
		return
	}

	// just one - default config
	if c.Server.TLS.CertFile != "" && c.Server.TLS.KeyFile != "" {
		certificate, certificateErr := tls.LoadX509KeyPair(c.Server.TLS.CertFile, c.Server.TLS.KeyFile)
		if certificateErr != nil {
			err = certificateErr
			return
		}
		certificates = append(certificates, certificate)
		return
	}

	// multiple certs for SNI - let us loop
	certificates = make([]tls.Certificate, len(c.Server.TLS.Certificates))
	for i, certConf := range c.Server.TLS.Certificates {
		certificate, certificateErr := tls.LoadX509KeyPair(certConf.CertFile, certConf.KeyFile)
		if certificateErr != nil {
			err = certificateErr
			return
		}
		certificates[i] = certificate
	}
	return
}

func (c *Config) hasCertificates() bool {
	return (c.Server.TLS.CertFile != "" && c.Server.TLS.KeyFile != "") || len(c.Server.TLS.Certificates) > 0
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
