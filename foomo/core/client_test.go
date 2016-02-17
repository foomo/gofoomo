package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/foomo/gofoomo/foomo"
	"github.com/foomo/gofoomo/foomo/bert"
)

type CoreConfig struct {
	EnabledModules   []string
	AvailableModules []string
	RootHttp         string
	buildNumber      int64
}

var testFoomo *foomo.Foomo

func poe(err error, msg string) {
	if err != nil {
		panic(msg + " : " + err.Error())
	}
}

func getTestServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			    "enabledModules": [
			        "Foomo"
			    ],
			    "availableModules": [
			        "Foomo"
			    ]
			}`))
	}))
	return ts
}

func getTestFoomo() *foomo.Foomo {

	if testFoomo == nil {
		ts := getTestServer()
		tmp := "/tmp" //os.TempDir()
		dir, err := ioutil.TempDir(tmp, "dummy-foomo")
		poe(err, "failed to get temp dir")
		bareFoomo, err := foomo.BareFoomo(dir, "test")
		poe(err, "failed to get bare foomo")
		b := bert.NewBert(bareFoomo)
		b.Prepare()
		f, err := foomo.NewFoomo(dir, "test", fmt.Sprint(ts.URL))
		if err != nil {
			panic("invalid test foomo " + err.Error())
		}
		testFoomo = f
	}
	return testFoomo
}

func TestGet(t *testing.T) {
	f := getTestFoomo()
	data, err := get(f, "config", "Foomo", "Foomo.core")
	if err != nil {
		t.Fatal(err)
	}
	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetJSON(t *testing.T) {
	f := getTestFoomo()
	config := new(CoreConfig)
	err := GetJSON(f, config, "config", "Foomo", "Foomo.core")
	if err != nil {
		t.Fatal(err)
	}
	if len(config.EnabledModules) < 1 {
		t.Fatal("there must be at least Foomo enabled")
	}
}

func TestGetConfig(t *testing.T) {
	f := getTestFoomo()
	config := new(CoreConfig)
	err := GetConfig(f, config, "Foomo", "Foomo.core", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(config.EnabledModules) < 1 {
		t.Fatal("there must be at least Foomo enabled")
	}
}
