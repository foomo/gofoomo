package bert

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func readConfig(filename string) (config map[string][]string, err error) {
	config = make(map[string][]string)
	yamlBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlBytes, config)
	return
}
