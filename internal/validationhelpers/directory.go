package validationhelpers

import (
	"errors"
	"os"
)

func CheckDirReadable(srcDir interface{}) error {
	assertedSrcDir, ok := srcDir.(string)
	if !ok {
		return errors.New("must be a string")
	}
	fileInfo, err := os.Stat(assertedSrcDir)
	if err != nil {
		return errors.New("must be accessible")
	}
	if !fileInfo.IsDir() {
		return errors.New("must be a directory path")
	}
	_, err = os.ReadDir(assertedSrcDir)
	if err != nil {
		return errors.New("must be readable")
	}
	return nil
}
