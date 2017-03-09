package proxy_test

import (
	"crypto/tls"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/foomo/gofoomo/proxy"
	"github.com/foomo/gofoomo/proxy/handler"
)

func getMockFile(name string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename) + string(os.PathSeparator) + "mock" + string(os.PathSeparator) + name
}

func TestConfig(t *testing.T) {
	// some very basic testing here ;)

	config, err := proxy.ReadConfig(getMockFile("config.yml"))
	if err != nil {
		t.Fatal("read error")
	}
	if config.Server.TLS.CertFile != "/foo.cert" {
		t.Fatal("could not find sever tls cert", config.Server)
	}
}

func TestSNI(t *testing.T) {
	config, err := proxy.ReadConfig(getMockFile("config-sni.yml"))
	config.Foomo.Dir = getMockFile(config.Foomo.Dir)
	if err != nil {
		t.Fatal("read error")
	}
	for i, certConf := range config.Server.TLS.Certificates {
		config.Server.TLS.Certificates[i].CertFile = getMockFile(certConf.CertFile)
		config.Server.TLS.Certificates[i].KeyFile = getMockFile(certConf.KeyFile)
	}
	p, err := proxy.NewServer(config)
	if err != nil {
		t.Fatal("could not instatiate proxy", err)
	}
	p.Proxy.AddHandler(handler.NewStaticFiles(p.Foomo))
	go func() {
		serveErr := p.ListenAndServe()
		if serveErr != nil {
			t.Fatal("could not listen and serve")
		}
	}()
	time.Sleep(time.Millisecond * 200)
	c := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	response, responseErr := c.Get("https://localhost:8443/foomo/modulesVar/Foomo.JS/test.js")
	if responseErr != nil {
		t.Fatal("failed to get", responseErr)
	}
	if response.TLS.PeerCertificates[0].Subject.CommonName != "localhost" {
		t.Fatal("SNI Fail, unexpected common name in first peer certificate common name:", response.TLS.PeerCertificates[0].Subject.CommonName)
	}
}
