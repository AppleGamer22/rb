package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func Source2TargetPath(sourceFilePath, sourcePathRoot, targetPathRoot string) (string, error) {
	relativePath, err := filepath.Rel(sourcePathRoot, sourceFilePath)
	if err != nil {
		return "", err
	}
	targetFilePath, err := filepath.Abs(fmt.Sprintf("%s/%s", targetPathRoot, relativePath))
	if err != nil {
		return "", err
	}
	return targetFilePath, nil
}

func DoesTargetFileExist(sourceFilePath, sourcePathRoot, targetPathRoot string) (bool, error) {
	targetFilePath, err := Source2TargetPath(sourceFilePath, sourcePathRoot, targetPathRoot)
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(targetFilePath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func WaitForDirectory(path string) {
	var searching = true
	for searching {
		_, err := os.Stat(path)
		if err != nil {
			fmt.Printf("Waiting for directory %s to be available...\n", path)
			time.Sleep(2 * time.Second)
		} else {
			searching = false
		}
	}
}