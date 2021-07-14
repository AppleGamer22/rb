package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

func BackupLog(fileNum, numberOfFiles int, file FileMetadata) {
	fmt.Printf("file #%d / %d (%s -> %s)\n", fileNum, numberOfFiles, file.SourcePath, file.TargetPath)
}

func Backup(filePath, targetPath string) error {
	var filesLog, err = GetLogFromFile(filePath)
	if err != nil {
		return err
	}

	for i, file := range filesLog.Files {
		BackupLog(i+1, len(filesLog.Files), file)
		err = Copy(file, targetPath)
		if err != nil {
			return err
		}
		filesLog.Files[i].Done = true
		err = SaveMetadataToFile(filesLog.Files, filePath, i+1, i == len(filesLog.Files)-1)
		if err != nil {
			return err
		}
		fmt.Println("done")
	}
	return nil
}

func Copy(file FileMetadata, targetPath string) error {
	fileStatSource, err := os.Stat(file.SourcePath)
	if err != nil {
		WaitForDirectory(targetPath)
	}
	if !fileStatSource.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", file.SourcePath)
	}
	_, err = os.Stat(file.TargetPath)
	if err != nil {
		err := os.MkdirAll(filepath.Dir(file.TargetPath), 0755)
		if err != nil {
			return err
		}
	}
	src, err := os.Open(file.SourcePath)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(file.TargetPath)

	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)

	return err
}

func SaveMetadataToFile(files []FileMetadata, path string, filesCopied int, keepTime bool) error {
	var dataLog = Log{Files: files, FilesCopied: filesCopied}
	if keepTime {
		dataLog.LastBackupTime = time.Now()
	} else {
		dataLog.LastBackupTime = time.Unix(0, 0)
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
