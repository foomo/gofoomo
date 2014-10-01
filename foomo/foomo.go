package foomo

import (
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	u "net/url"
	"strings"
)

const shaPrefix = "{SHA}"

type Foomo struct {
	Root                 string
	RunMode              string
	URL                  *u.URL
	basicAuthCredentials struct {
		user     string
		password string
	}
}

func NewFoomo(foomoDir string, runMode string, url string) (f *Foomo, err error) {
	f, err = makeFoomo(foomoDir, runMode, url, true)
	return
}

func makeFoomo(foomoDir string, runMode string, url string, init bool) (foomo *Foomo, err error) {
	f := new(Foomo)
	f.Root = foomoDir
	f.URL, err = u.Parse(url)
	f.RunMode = runMode
	if init {
		f.setupBasicAuthCredentials()
	}
	return f, err
}

func (f *Foomo) BasicAuth(domain string, user string, password string) bool {
	for _, line := range strings.Split(f.getBasicAuthFileContentsForDomain(domain), "\n") {
		lineParts := strings.Split(line, ":")
		if len(lineParts) == 2 && lineParts[0] == user {
			hash := getBasicAuthHash(password)
			return hash == lineParts[1]
		}
	}
	return false
}

func (f *Foomo) getBasicAuthFileContentsForDomain(domain string) string {
	basicAuthFilename := f.GetBasicAuthFilename(domain)
	bytes, err := ioutil.ReadFile(basicAuthFilename)
	if err != nil {
		return ""
	} else {
		return string(bytes)
	}
}

func (f *Foomo) setupBasicAuthCredentials() error {
	f.basicAuthCredentials.user = "gofoomo"
	f.basicAuthCredentials.password = makeToken(50)
	return ioutil.WriteFile(f.GetBasicAuthFilename("default"), []byte(setBasicAuthForUserInBasicAuthFileContents(f.getBasicAuthFileContentsForDomain("default"), f.basicAuthCredentials.user, f.basicAuthCredentials.password)), 0644)
}

func setBasicAuthForUserInBasicAuthFileContents(basicAuthFileContents string, user string, password string) string {
	newLines := make([]string, 0)
LineLoop:
	for _, line := range strings.Split(basicAuthFileContents, "\n") {
		lineParts := strings.Split(line, ":")
		if len(lineParts) == 2 && lineParts[0] == user {
			continue LineLoop
		} else if len(line) > 0 {
			newLines = append(newLines, line)
		}
	}
	newLines = append(newLines, user+":"+getBasicAuthHash(password))
	return strings.Join(newLines, "\n")
}

func getBasicAuthHash(password string) string {
	s := sha1.New()
	s.Write([]byte(password))
	passwordSum := []byte(s.Sum(nil))
	return shaPrefix + base64.StdEncoding.EncodeToString(passwordSum)
}

func (f *Foomo) GetURLWithCredentialsForDefaultBasicAuthDomain() string {
	url, _ := u.Parse(f.URL.String())
	url.User = u.UserPassword(f.basicAuthCredentials.user, f.basicAuthCredentials.password)
	return url.String()
}

func (f *Foomo) GetBasicAuthCredentialsForDefaultBasicAuthDomain() (user string, password string) {
	return f.basicAuthCredentials.user, f.basicAuthCredentials.password
}

func (f *Foomo) GetModuleDir(moduleName string, dir string) string {
	return f.Root + "/modules/" + moduleName + "/" + dir
}

func (f *Foomo) GetVarDir() string {
	return f.Root + "/var/" + f.RunMode
}

func (f *Foomo) GetModuleHtdocsDir(moduleName string) string {
	return f.GetModuleDir(moduleName, "htdocs")
}

func (f *Foomo) GetModuleHtdocsVarDir(moduleName string) string {
	return f.GetVarDir() + "/htdocs/modulesVar/" + moduleName
}

func (f *Foomo) GetBasicAuthFilename(domain string) string {
	return f.GetVarDir() + "/basicAuth/" + domain
}
