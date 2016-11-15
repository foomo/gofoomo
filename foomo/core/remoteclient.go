package core

import "net/url"

type RemoteClient struct {
	url string
}

func NewRemoteClient(urlString string) (rc *RemoteClient, err error) {
	_, e := url.Parse(urlString)
	if e != nil {
		err = e
		return
	}
	rc = &RemoteClient{
		url: urlString,
	}
	return
}

func (rc *RemoteClient) get(path ...string) (data []byte, err error) {
	return get(rc.url, path...)
}

func (rc *RemoteClient) GetConfig(target interface{}, moduleName string, configName string, domain string) (err error) {
	if len(domain) == 0 {
		return getJSON(rc.url, target, "config", moduleName, configName)
	}
	return getJSON(rc.url, target, "config", moduleName, configName, domain)
}
