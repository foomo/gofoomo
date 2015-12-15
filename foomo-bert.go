package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/foomo/gofoomo/foomo"
	"github.com/foomo/gofoomo/foomo/bert"
	"github.com/foomo/gofoomo/utils"
)

type foomoFlagsPrepare struct {
	runMode *string
	dir     *string
}

type foomoFlagsReset struct {
	foomoFlagsPrepare
	address    *string
	mainModule *string
}

func getFlagRunMode(fs *flag.FlagSet) *string {
	return fs.String("run-mode", "", "foomo run mode test | development | production")
}

func getFlagDir(fs *flag.FlagSet) *string {
	return fs.String("dir", "", "path/to/your/foomo/root")
}

func foomoFlagsetReset() (fs *flag.FlagSet, f *foomoFlagsReset) {
	f = &foomoFlagsReset{}
	fs = flag.NewFlagSet(os.Args[0]+" foomo-prepare", flag.ContinueOnError)
	f.runMode = getFlagRunMode(fs)
	f.dir = getFlagDir(fs)
	f.address = fs.String("addr", "", "address of the foomo server")
	f.mainModule = fs.String("main-module", "Foomo", "name of main module")
	return fs, f
}

func foomoFlagsetPrepare() (fs *flag.FlagSet, f *foomoFlagsPrepare) {
	f = &foomoFlagsPrepare{}
	fs = flag.NewFlagSet(os.Args[0]+" prepare", flag.ContinueOnError)
	f.runMode = getFlagRunMode(fs)
	f.dir = getFlagDir(fs)
	return fs, f
}

func validateFlagsReset(f *foomoFlagsReset) (err error) {
	fp := &foomoFlagsPrepare{
		runMode: f.runMode,
		dir:     f.dir,
	}
	prepareErr := validateFlagsPrepare(fp)
	if prepareErr != nil {
		return prepareErr
	}
	// main module
	if len(*f.mainModule) == 0 {
		return errors.New("main module must be set missing -main-module")
	}
	moduleDir := path.Join(*f.dir, "modules", *f.mainModule)
	_, dirErr := utils.IsDir(moduleDir)
	if dirErr != nil {
		return errors.New("main module dir error: " + dirErr.Error())
	}
	// addr
	if len(*f.address) == 0 {
		return errors.New("missing address -addr")
	}
	_, err = url.Parse(*f.address)
	if err != nil {
		return errors.New(fmt.Sprint("could not parse address"))
	}
	return nil
}

func validateFlagsPrepare(f *foomoFlagsPrepare) (err error) {
	// run mode
	switch *f.runMode {
	case foomo.RunModeTest, foomo.RunModeProduction, foomo.RunModeDevelopment:
	default:
		return errors.New(fmt.Sprintln("invalid run mode", "\""+*f.runMode+"\"", "must be one of", []string{foomo.RunModeTest, foomo.RunModeProduction, foomo.RunModeDevelopment}))
	}
	// foomo dir
	_, dirErr := utils.IsDir(*f.dir)
	if err != nil {
		return errors.New("dir is not a directory: " + dirErr.Error())
	}
	return nil
}

func usage() {
	fmt.Println("usage:", os.Args[0], "<command>")
	fsPrepare, _ := foomoFlagsetPrepare()
	fsReset, _ := foomoFlagsetReset()
	for command, fs := range map[string]*flag.FlagSet{
		"prepare": fsPrepare,
		"reset":   fsReset,
	} {
		fmt.Println(os.Args[0], command, ":")
		fs.PrintDefaults()
	}
}

func flagErr(fs *flag.FlagSet, comment string, err error) {
	if err != nil {
		fmt.Println(comment, err.Error())
		usage()
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "prepare":
			fs, flagsPrepare := foomoFlagsetPrepare()
			fs.Parse(os.Args[2:])
			err := validateFlagsPrepare(flagsPrepare)
			flagErr(fs, "prepare", err)
			fmt.Println("preparing foomo in", *flagsPrepare.dir, "to run in run mode", *flagsPrepare.runMode)
			f, foomoErr := foomo.NewFoomo(*flagsPrepare.dir, *flagsPrepare.runMode, "fake://no-where:8080")
			if foomoErr != nil {
				fmt.Println(foomoErr.Error())
				os.Exit(1)
			}
			b := bert.NewBert(f)
			prepareErr := b.Prepare()
			if prepareErr != nil {
				fmt.Println("failed to prepare", prepareErr.Error())
				os.Exit(1)
			}
		case "reset":
			fs, flagsReset := foomoFlagsetReset()
			fs.Parse(os.Args[2:])
			err := validateFlagsReset(flagsReset)
			flagErr(fs, "reset", err)
			fmt.Println("resetting foomo in", *flagsReset.dir, "in run mode", *flagsReset.runMode)
			f, foomoErr := foomo.NewFoomo(*flagsReset.dir, *flagsReset.runMode, *flagsReset.address)
			if foomoErr != nil {
				fmt.Println(foomoErr.Error())
				os.Exit(1)
			}
			b := bert.NewBert(f)
			resetErr := b.Reset(*flagsReset.mainModule)
			if resetErr != nil {
				fmt.Println("failed to reset", resetErr.Error())
				os.Exit(1)
			}
		default:
			usage()
		}
	} else {
		usage()
	}
}
