package bert

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/bgentry/speakeasy"
	"github.com/foomo/gofoomo/foomo"
	"github.com/foomo/htpasswd"
)

// Prepare an installation (very basic skeleton)
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

// PrepareAdmin asks and adds basic auth credentials for foomo default auth\
// domain
func (b *Bert) PrepareAdmin(admin string) error {

	fmt.Println("adding admin user", admin)

	passwd, err := speakeasy.Ask("enter password for " + admin + " 🔑 :\n")
	if err != nil {
		return errors.New("could not read password, giving up")
	}
	passwordFile := b.foomo.GetBasicAuthFilename(foomo.DefaultBasicAuthDomainName)
	fmt.Println("adding password for admin in:", passwordFile)

	err = htpasswd.SetPassword(passwordFile, admin, passwd, htpasswd.HashBCrypt)
	if err != nil {
		return errors.New("could not write default basic auth password for admin: " + err.Error())
	}
	fmt.Println("added password")
	fmt.Println("DONE")
	return nil
}
