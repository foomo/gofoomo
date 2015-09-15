package proxy

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v1"
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
	TLSModeStrict  string = "strict"
	TLSModeLoose          = "loose"
	TLSModeDefault        = "default"
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
