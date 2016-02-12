package proxy

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Auth struct {
	Domain string
	Realm  string
}

type TLS struct {
	Mode     string
	Address  string
	CertFile string
	KeyFile  string
}

const (
	// this is serious and we do not mind loosing clients (= Mozilla "modern" compatibility)
	// Compatible clients have versions equal or greater than Firefox 27, Chrome 22, IE 11, Opera 14, Safari 7, Android 4.4, Java 8
	TLSModeStrict = "strict"
	// ecommerce compromise
	// Compatible clients (>=): Firefox 1, Chrome 1, IE 7, Opera 5, Safari 1, Windows XP IE8, Android 2.3, Java 7
	TLSModeLoose = "loose"
	// standard crypto/tls.Config un touched - highly compatible, but possibly insecure
	TLSModeDefault = "default"
)

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
		config.Server.TLS.Mode = TLSModeDefault
	}
	switch config.Server.TLS.Mode {
	case TLSModeDefault, TLSModeLoose, TLSModeStrict:
	default:
		return nil, errors.New("unknown server.tls.mode: " + config.Server.TLS.Mode + " - must be one of: " + fmt.Sprintln([]string{TLSModeDefault, TLSModeLoose, TLSModeStrict}))
	}
	return config, nil
}
