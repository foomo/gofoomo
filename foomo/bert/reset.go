package bert

import "errors"

// Reset reset a bert installation
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
