package proxy

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
)

type Auth struct {
	Domain string
	Realm  string
}

type Config struct {
	// how should the proxy server run
	Server struct {
		Address string
		Auth    *Auth
		TLS     struct {
			Address  string
			CertFile string
			KeyFile  string
		}
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
	return config, nil
}
