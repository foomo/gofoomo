package utils

import (
	"errors"
	"os"
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
