package foomo

import (
	"errors"
	"fmt"
	u "net/url"
	"os"

	"github.com/foomo/htpasswd"
)

const (
	// RunModeTest run mode test
	RunModeTest = "test"
	// RunModeDevelopment run mode development
	RunModeDevelopment = "development"
	// RunModeProduction run mode production
	RunModeProduction = "production"
	// DefaultBasicAuthDomainName name of default basic auth domain
	DefaultBasicAuthDomainName = "default"
)

// Foomo foomo go wrapper
type Foomo struct {
	Root                 string
	RunMode              string
	URL                  *u.URL
	basicAuthCredentials struct {
		user     string
		password string
	}
}

// NewFoomo get a foomo instance
func NewFoomo(foomoDir string, runMode string, url string) (f *Foomo, err error) {
	return makeFoomo(foomoDir, runMode, url, true)
}

// BareFoomo is an unistalled foomo, that bert uses to prepare an installation
func BareFoomo(foomoDir string, runMode string) (f *Foomo, err error) {
	return makeFoomo(foomoDir, runMode, "fake://no", false)
}

func makeFoomo(foomoDir string, runMode string, address string, init bool) (f *Foomo, err error) {
	// validate run mode
	switch runMode {
	case RunModeTest, RunModeDevelopment, RunModeProduction:
	default:
		return nil, errors.New("invalid run mode: " + runMode + " must be one of: " + fmt.Sprintln([]string{RunModeTest, RunModeDevelopment, RunModeProduction}))
	}
	// validate root dir
	_, err = os.Stat(foomoDir)
	if err != nil {
		return nil, errors.New("can not access foomo dir: " + err.Error())
	}

	// validate url
	if len(address) == 0 {
		return nil, errors.New("foomo address must not be empty")
	}
	u, err := u.Parse(address)
	if err != nil {
		return nil, errors.New("can not parse foomo url: " + err.Error())
	}
	// instantiate
	f = &Foomo{
		RunMode: runMode,
		URL:     u,
		Root:    foomoDir,
	}
	// init
	if init {
		authErr := f.setupBasicAuthCredentials()
		if authErr != nil {
			return nil, errors.New("can not set up auth: " + authErr.Error())
		}
	}
	return f, err
}

func (f *Foomo) setupBasicAuthCredentials() error {
	if f.URL.User != nil {
		f.basicAuthCredentials.user = f.URL.User.Username()
		password, passwordOK := f.URL.User.Password()
		if passwordOK {
			f.basicAuthCredentials.password = password
			return nil
		}
	}
	f.basicAuthCredentials.user = "gofoomo"
	f.basicAuthCredentials.password = makeToken(50)
	return htpasswd.SetPassword(f.GetBasicAuthFilename("default"), f.basicAuthCredentials.user, f.basicAuthCredentials.password, htpasswd.HashBCrypt)
}

// GetURLWithCredentialsForDefaultBasicAuthDomain i.e. sth. like http(s)://user:password@foomo-server.org(:8080)
func (f *Foomo) GetURLWithCredentialsForDefaultBasicAuthDomain() string {
	url, _ := u.Parse(f.URL.String())
	url.User = u.UserPassword(f.basicAuthCredentials.user, f.basicAuthCredentials.password)
	return url.String()
}

// GetBasicAuthCredentialsForDefaultBasicAuthDomain user, password generated for the local foomo instance
func (f *Foomo) GetBasicAuthCredentialsForDefaultBasicAuthDomain() (user string, password string) {
	return f.basicAuthCredentials.user, f.basicAuthCredentials.password
}

// GetModuleDir root dir of a module
func (f *Foomo) GetModuleDir(moduleName string, dir string) string {
	return f.Root + "/modules/" + moduleName + "/" + dir
}

// GetVarDir root var dir for the current run mode
func (f *Foomo) GetVarDir() string {
	return f.Root + "/var/" + f.RunMode
}

// GetModuleHtdocsDir htdocs dir in a module
func (f *Foomo) GetModuleHtdocsDir(moduleName string) string {
	return f.GetModuleDir(moduleName, "htdocs")
}

// GetModuleCacheDir cache dir for a module
func (f *Foomo) GetModuleCacheDir(moduleName string) string {
	return f.GetVarDir() + "/cache/" + moduleName
}

// GetModuleHtdocsVarDir module htdocs var dir
func (f *Foomo) GetModuleHtdocsVarDir(moduleName string) string {
	return f.GetVarDir() + "/htdocs/modulesVar/" + moduleName
}

// GetBasicAuthFilename basic auth file name for a domain
func (f *Foomo) GetBasicAuthFilename(domain string) string {
	return f.GetVarDir() + "/basicAuth/" + domain
}
