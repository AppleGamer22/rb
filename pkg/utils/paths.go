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