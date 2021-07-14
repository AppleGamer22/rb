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

func BackupLog(fileNum, numberOfFiles int, file FileMetadata) {
	fmt.Printf("file #%d / %d (%s -> %s)\n", fileNum, numberOfFiles,
		file.SourcePath, file.TargetPath)
}

func Backup(filePath string) error {
	var data, err = ioutil.ReadFile(filePath)

	if err != nil {
		return err
	}

	var filesLog Log
	err = json.Unmarshal(data, &filesLog)
	if err != nil {
		return err
	}

	for i, file := range filesLog.Files {
		err = Copy(file)
		if err != nil {
			return err
		}
		filesLog.Files[i].Done = true
		err = SaveMetadataToFile(filesLog.Files, filePath, i + 1, i == len(filesLog.Files) - 1)
		if err != nil {
			return err
		}
		BackupLog(i+1, len(filesLog.Files), file)
	}
	return nil
}

func Copy(file FileMetadata) error {
	fileStatSource, err := os.Stat(file.SourcePath)
	if err != nil {
		return err
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

// func MarkAsDone(path string, i int) error {
// 	var data, err1 = ioutil.ReadFile(path)
// 	if err1 != nil {
// 		return err1
// 	}
// 	var metadata []FileMetadata
// 	var err2 = json.Unmarshal(data, &metadata)
// 	if err2 != nil {
// 		return err2
// 	}
// 	if len(metadata) > 0 && 0 <= i && i < len(metadata) {
// 		metadata[i].Done = true
// 	} else {
// 		return errors.New("index is out of scope")
// 	}
// 	var err3 = SaveMetadataToFile(metadata, path, i, )
// 	if err3 != nil {
// 		return err3
// 	}
// 	return nil
// }
