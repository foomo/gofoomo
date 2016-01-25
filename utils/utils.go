package utils

import (
	"errors"
	"os"
	"path"
	"runtime"
)

// IsDir tells you is it a dir or not
func IsDir(filename string) (isDir bool, err error) {
	i, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	if !i.IsDir() {
		return false, errors.New("not a dir")
	}
	return true, nil
}

// __DIR__
func GetCurrentDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
