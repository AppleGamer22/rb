package rb

import (
	"fmt"
	"github.com/AppleGamer22/recursive-backup/internal/utils"
	"io"
	"os"
	"path/filepath"
	"time"
)

func printProgressMessage(fileNum int, sourcePath, targetPath string) {
	fmt.Printf("file #%d (%s -> %s)\n", fileNum, sourcePath, targetPath)
}

// Backs up file file to its correct location on target based on its target sub-direcory
func (rber RecursiveBackupper) backupFile(sourceFilePath string, fileNumber int) (targetFilePath string, copyTime time.Time, err error) {
	targetFilePath, err = utils.Source2TargetPath(sourceFilePath, rber.SourceRoot, rber.TargetRoot)
	if err != nil {
		return "", time.Unix(0, 0), err
	}
	printProgressMessage(fileNumber, sourceFilePath, targetFilePath)
	copyTime, err = rber.copyFile(sourceFilePath, targetFilePath)
	if err != nil {
		return "", time.Unix(0, 0), err
	}
	fmt.Println("done")
	return targetFilePath, copyTime, nil
}

// Copies file to the provided destination
func (rber RecursiveBackupper) copyFile(sourcePath, targetPath string) (time.Time, error) {
	fileStatSource, err := os.Stat(sourcePath)
	if err != nil {
		utils.WaitForDirectory(targetPath)
		return rber.copyFile(sourcePath, targetPath)
	}
	if !fileStatSource.Mode().IsRegular() {
		return time.Unix(0, 0), fmt.Errorf("%s is not a regular file", sourcePath)
	}
	_, err = os.Stat(targetPath)
	if err != nil {
		err := os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			utils.WaitForDirectory(rber.TargetRoot)
			return rber.copyFile(sourcePath, targetPath)
		}
	}

	src, err := os.Open(sourcePath)
	if err != nil {
		return time.Unix(0, 0), err
	}
	defer src.Close()

	dest, err := os.Create(targetPath)
	if err != nil {
		utils.WaitForDirectory(rber.TargetRoot)
		return rber.copyFile(sourcePath, targetPath)
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	return time.Now(), err
}
