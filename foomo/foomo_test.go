package foomo

import (
	"io/ioutil"
	"os"
	"testing"
)

var tmp string

func getTestFoomoForFSStuff() *Foomo {
	tempDir, err := ioutil.TempDir(os.TempDir(), "dummy-foomo")
	if err != nil {
		panic(err)
	}
	tempDir = tempDir[0 : len(tempDir)-1]
	tmp = tempDir
	os.MkdirAll(tempDir, 0777)
	f, err := BareFoomo(tempDir, "test")
	if err != nil {
		panic(err)
	}

	return f
}

func assertTempPath(t *testing.T, topic string, expected string, actual string) {
	assertStringsEqual(t, topic, tmp+"/"+expected, actual)
}
func assertStringsEqual(t *testing.T, topic string, expected string, actual string) {
	if actual != expected {
		t.Fatal(topic, "actual: ", actual, " != expected: ", expected)
	}
}

func TestGetVarDir(t *testing.T) {
	actual := getTestFoomoForFSStuff().GetVarDir()
	expected := "var/test"
	assertTempPath(t, "var dir", expected, actual)
}

func TestGetModuleDir(t *testing.T) {
	assertTempPath(
		t,
		"module dir",
		"modules/Foomo/htdocs",
		getTestFoomoForFSStuff().GetModuleDir("Foomo", "htdocs"),
	)
}

func TestGetModuleHtdocsDir(t *testing.T) {
	assertTempPath(
		t,
		"module htdocs dir",
		"modules/Foomo/htdocs",
		getTestFoomoForFSStuff().GetModuleHtdocsDir("Foomo"),
	)
}

func TestGetModuleHtdocsVarDir(t *testing.T) {
	assertTempPath(
		t,
		"module htdocs var dir",
		"var/test/htdocs/modulesVar/Foomo",
		getTestFoomoForFSStuff().GetModuleHtdocsVarDir("Foomo"),
	)
}

func TestGetBasicAuthFilename(t *testing.T) {
	assertTempPath(
		t,
		"basic auth file",
		"var/test/basicAuth/sepp",
		getTestFoomoForFSStuff().GetBasicAuthFilename("sepp"),
	)
}
