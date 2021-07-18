package rb

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/AppleGamer22/recursive-backup/pkg/utils"
)

func printProgressMessage(fileNum int, sourcePath, targetPath string) {
	fmt.Printf("file #%d (%s -> %s)\n", fileNum, sourcePath, targetPath)
}

func BackupFile(sourceFilePath, logPath, sourcePathRoot, targetPathRoot string, i int, startTime time.Time) (string, string, time.Time, error) {
	targetFilePath, err := utils.Source2TargetPath(sourceFilePath, sourcePathRoot, targetPathRoot)
	if err != nil {
		return "", "", time.Unix(0, 0), err
	}
	printProgressMessage(i, sourceFilePath, targetFilePath)
	copyTime, err := CopyFile(sourceFilePath, targetFilePath, targetPathRoot)
	if err != nil {
		return "", "", time.Unix(0, 0), err
	}
	return sourceFilePath, targetFilePath, copyTime, nil
}

func CopyFile(sourcePath, targetPath string, targetPathRoot string) (time.Time, error) {
	fileStatSource, err := os.Stat(sourcePath)
	if err != nil {
		utils.WaitForDirectory(targetPathRoot)
		return CopyFile(sourcePath, targetPath, targetPathRoot)
	}
	if !fileStatSource.Mode().IsRegular() {
		return time.Unix(0, 0), fmt.Errorf("%s is not a regular file", sourcePath)
	}
	_, err = os.Stat(targetPath)
	if err != nil {
		err := os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			utils.WaitForDirectory(targetPathRoot)
			return CopyFile(sourcePath, targetPath, targetPathRoot)
		}
	}

	src, err := os.Open(sourcePath)
	if err != nil {
		return time.Unix(0, 0), err
	}
	defer src.Close()

	dest, err := os.Create(targetPath)
	if err != nil {
		utils.WaitForDirectory(targetPathRoot)
		return CopyFile(sourcePath, targetPath, targetPathRoot)
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	return time.Now(), err
}


