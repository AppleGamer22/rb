package rb

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetLogFromFile(path string) (Log, error) {
	var data, err = ioutil.ReadFile(path)

	if err != nil {
		return Log{}, err
	}

	var filesLog Log
	err = json.Unmarshal(data, &filesLog)
	if err != nil {
		return Log{}, err
	}
	return filesLog, nil
}

func BackupLog(fileNum, numberOfFiles int, sourcePath, targetPath string) {
	fmt.Printf("file #%d / %d (%s -> %s)\n", fileNum, numberOfFiles, sourcePath, targetPath)
}

func Backup(sourcesLogPath, sourcePathRoot, targetPathRoot string, fileCount int, startTime time.Time) error {
	fileSourcesLog, err := os.Open(sourcesLogPath)
	if err != nil {
		return err
	}
	defer fileSourcesLog.Close()
	var reader = bufio.NewReader(fileSourcesLog)
	var targetsLogPath = fmt.Sprintf("rb_target_%d-%d-%d_%d:%d:%d.csv", startTime.Day(), startTime.Month(), startTime.Year(), startTime.Hour(), startTime.Minute(), startTime.Second())
	fileTargetsLog, err := os.Create(targetsLogPath)
	if err != nil {
		return err
	}
	defer fileTargetsLog.Close()
	var writer = bufio.NewWriter(fileTargetsLog)
	for i := 1;; i++ {
		var data, err1 = reader.ReadString('\n')
		if err1 == io.EOF {
			break
		}
		var sourcePath = strings.Replace(string(data), "\n", "", -1)
		var relativePath, err2 = filepath.Rel(sourcePathRoot, sourcePath)
		if err2 != nil {
			return err2
		}
		var targetPath, err3 = filepath.Abs(fmt.Sprintf("%s/%s", targetPathRoot, relativePath))
		if err3 != nil {
			return err3
		}
		BackupLog(i, fileCount, sourcePath, targetPath)
		var modTime, err4 = Copy(sourcePath, targetPath, targetPathRoot)
		if err4 != nil {
			return err4
		}
		writer.WriteString(fmt.Sprintf("%s,%s,%s\n", sourcePath, targetPath, modTime))
		fmt.Println("done")
	}
	return nil
}

func Copy(sourcePath, targetPath string, targetPathRoot string) (time.Time, error) {
	fileStatSource, err := os.Stat(sourcePath)
	if err != nil {
		WaitForDirectory(targetPathRoot)
		return Copy(sourcePath, targetPath, targetPathRoot)
	}
	if !fileStatSource.Mode().IsRegular() {
		return time.Unix(0, 0), fmt.Errorf("%s is not a regular file", sourcePath)
	}
	_, err = os.Stat(targetPath)
	if err != nil {
		err := os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			WaitForDirectory(targetPathRoot)
			return Copy(sourcePath, targetPath, targetPathRoot)
		}
	}
	src, err := os.Open(sourcePath)
	if err != nil {
		return time.Unix(0, 0), err
	}
	defer src.Close()

	dest, err := os.Create(targetPath)

	if err != nil {
		WaitForDirectory(targetPathRoot)
		return Copy(sourcePath, targetPath, targetPathRoot)
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	return fileStatSource.ModTime(), err
}

func SaveMetadataToFile(files []FileMetadata, path string, filesCopied int, keepTime bool, date time.Time) error {
	var dataLog = Log{Files: files, FilesCopied: filesCopied}
	if keepTime {
		dataLog.LastBackupTime = time.Now()
	} else {
		dataLog.LastBackupTime = date
	}
	var json, err1 = json.MarshalIndent(dataLog, "", "\t")
	if err1 != nil {
		return err1
	}
	var err2 = ioutil.WriteFile(path, json, 0644)
	if err2 != nil {
		return err2
	}
	return nil
}

func WaitForDirectory(path string) {
	fmt.Printf("Waiting for directory %s to be available...\n", path)
	var searching = true
	for searching {
		_, err := os.Stat(path)
		if err != nil {
			time.Sleep(2 * time.Second)
		} else {
			searching = false
		}
	}
}
