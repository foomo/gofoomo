package foomo

import (
	"strings"
	"testing"
)

func TestSetBasicAuthForUserInBasicAuthFileContents(t *testing.T) {
	ba := "foo:bar\ntest:gone\nhansi:toll"
	newBa := setBasicAuthForUserInBasicAuthFileContents(ba, "test", "test")
	if len(strings.Split(newBa, "\n")) != 3 {
		t.Fatal("wrong line count")
	}
}

func getTestFoomoForFSStuff() *Foomo {
	f, _ := makeFoomo("/var/www/foomo", "test", "http://test.foomo", false)
	return f
}

func assertStringsEqual(t *testing.T, topic string, expected string, actual string) {
	if actual != expected {
		t.Fatal(topic, "actual: ", actual, " != expected: ", expected)
	}
}

func TestGetVarDir(t *testing.T) {
	actual := getTestFoomoForFSStuff().GetVarDir()
	expected := "/var/www/foomo/var/test"
	assertStringsEqual(t, "var dir", expected, actual)
}

func TestGetModuleDir(t *testing.T) {
	assertStringsEqual(t, "module dir", "/var/www/foomo/modules/Foomo/htdocs", getTestFoomoForFSStuff().GetModuleDir("Foomo", "htdocs"))
}

func TestGetModuleHtdocsDir(t *testing.T) {
	assertStringsEqual(t, "module htdocs dir", "/var/www/foomo/modules/Foomo/htdocs", getTestFoomoForFSStuff().GetModuleHtdocsDir("Foomo"))
}

func TestGetModuleHtdocsVarDir(t *testing.T) {
	assertStringsEqual(t, "module htdocs var dir", "/var/www/foomo/var/test/htdocs/modulesVar/Foomo", getTestFoomoForFSStuff().GetModuleHtdocsVarDir("Foomo"))
}

func TestGetBasicAuthFilename(t *testing.T) {
	assertStringsEqual(t, "basic auth file", "/var/www/foomo/var/test/basicAuth/sepp", getTestFoomoForFSStuff().GetBasicAuthFilename("sepp"))
}
