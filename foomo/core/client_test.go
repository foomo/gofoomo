package core

import (
	"encoding/json"
	"github.com/foomo/gofoomo/foomo"
	"testing"
)

type CoreConfig struct {
	EnabledModules   []string
	AvailableModules []string
	RootHttp         string
	buildNumber      int64
}

var testFoomo *foomo.Foomo

func getTestFoomo() *foomo.Foomo {
	if testFoomo == nil {
		f, _ := foomo.NewFoomo("/Users/jan/vagrant/schild/www/schild", "test", "http://schild-local-test.bestbytes.net")
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
