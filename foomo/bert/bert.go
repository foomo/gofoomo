package bert

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/foomo/gofoomo/foomo"
)

type Bert struct {
	foomo *foomo.Foomo
}

func NewBert(f *foomo.Foomo) *Bert {
	return &Bert{
		foomo: f,
	}
}

type moduleList struct {
	EnabledModules []string
}

const line = "-------------------------------------------------------------------------------"

func (b *Bert) Reset(mainModuleName string) error {
	err := b.get("resetting everything", "/foomo/hiccup.php?class=hiccup&action=resetEverything")
	if err != nil {
		return errors.New("failed to reset everything: " + err.Error())
	}

	err = b.get("enabling main module "+mainModuleName, "/foomo/core.php/enableModule/"+mainModuleName)
	if err != nil {
		return errors.New("enabling main module failed" + err.Error())
	}

	err = b.get("trying to create missing module resources", "/foomo/core.php/tryCreateModuleResources")
	if err != nil {
		return errors.New("failed to create module resources" + err.Error())
	}

	err = b.get("running make clean all", "/foomo/core.php/make/clean,all")
	if err != nil {
		return errors.New("make clean, all failed: " + err.Error())
	}
	return nil
}

func (b *Bert) get(explanation string, path string) error {
	fmt.Println(explanation)
	fmt.Println(line)
	url := b.foomo.GetURLWithCredentialsForDefaultBasicAuthDomain() + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println("something did not return a 200 for path:", path, ", status code:", response.StatusCode, ", status:", response.Status)
	}
	fmt.Println(string(responseBytes))
	fmt.Println(line)
	return nil
}

func (b *Bert) Prepare() error {
	b.foomo.GetVarDir()
	dirs := []string{
		"composer",
		"docker",
		path.Join("config", b.foomo.RunMode, "foomo"),
		path.Join("var", b.foomo.RunMode),
		path.Join("var", b.foomo.RunMode, "logs"),
		path.Join("var", b.foomo.RunMode, "basicAuth"),
		path.Join("var", b.foomo.RunMode, "sessions"),
		path.Join("var", b.foomo.RunMode, "tmp"),
		path.Join("var", b.foomo.RunMode, "logs"),
		path.Join("var", b.foomo.RunMode, "cache"),
		path.Join("var", b.foomo.RunMode, "modules"),
		path.Join("var", b.foomo.RunMode, "htdocs", "modules"),
		path.Join("var", b.foomo.RunMode, "htdocs", "modulesVar"),
	}
	for _, dir := range dirs {
		dirToMake := path.Join(b.foomo.Root, dir)
		mkdirErr := os.MkdirAll(dirToMake, 0744)
		if mkdirErr != nil {
			return mkdirErr
		}
		fmt.Println("created path", dirToMake)
	}
	return nil
}
