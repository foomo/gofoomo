package proxy

import (
	"os"
	"path"
	"runtime"
	"testing"
)

func getMockFile(name string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename) + string(os.PathSeparator) + "mock" + string(os.PathSeparator) + name
}

func TestConfig(t *testing.T) {
	// some very basic testing here ;)

	config, err := ReadConfig(getMockFile("config.yml"))
	if err != nil {
		t.Fatal("read error")
	}
	if config.Server.TLS.CertFile != "/foo.cert" {
		t.Fatal("could not find sever tls cert", config.Server)
	}
}
