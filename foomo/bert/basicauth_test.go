package bert

import (
	"path"
	"testing"

	"github.com/foomo/gofoomo/utils"
)

func TestReadConfig(t *testing.T) {
	config, err := readConfig(path.Join(utils.GetCurrentDir(), "basicauth_test", "test.yml"))
	utils.PanicOnErr(err)
	expectedAuthDomains := map[string][]string{
		"default": []string{
			"secret/foo",
			"secret/bar",
		},
		"sepp": []string{
			"secret/sepp/test",
		},
	}
	for expectedAuthDomain, secrets := range expectedAuthDomains {
		configSecrets, ok := config[expectedAuthDomain]
		if !ok {
			t.Fatal("missing auth domain", expectedAuthDomain)
		}
		for index, secret := range secrets {
			configSecret := configSecrets[index]
			if configSecret != secret {
				t.Fatal("unexpected secret", secret, "should be", configSecret)
			}
		}
	}
}
