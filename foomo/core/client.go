package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/foomo/gofoomo/foomo"
)

func get(foomo *foomo.Foomo, path ...string) (data []byte, err error) {
	callUrl := foomo.GetURLWithCredentialsForDefaultBasicAuthDomain()
	encodedPath := ""
	for _, pathEntry := range path {
		encodedPath += "/" + url.QueryEscape(pathEntry)
	}
	req, err := http.NewRequest("GET", callUrl+"/foomo/core.php"+encodedPath, nil)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unfriendly answer: " + resp.Status)
	}
	if err == nil {
		// handle error
		defer resp.Body.Close()
		data, err = ioutil.ReadAll(resp.Body)
	}
	return data, err
}

func GetJSON(foomo *foomo.Foomo, target interface{}, path ...string) error {
	data, err := get(foomo, path...)
	if err == nil {
		return json.Unmarshal(data, &target)
	} else {
		return err
	}
}

func GetConfig(foomo *foomo.Foomo, target interface{}, moduleName string, configName string, domain string) (err error) {
	if len(domain) == 0 {
		return GetJSON(foomo, target, "config", moduleName, configName)
	} else {
		return GetJSON(foomo, target, "config", moduleName, configName, domain)
	}
}
