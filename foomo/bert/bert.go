package bert

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/foomo/gofoomo/foomo"
)

// Bert is a foomo helper
type Bert struct {
	foomo *foomo.Foomo
}

// NewBert constructor
func NewBert(f *foomo.Foomo) *Bert {
	return &Bert{
		foomo: f,
	}
}

type moduleList struct {
	EnabledModules []string
}

const line = "-------------------------------------------------------------------------------"

// plain http get
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
