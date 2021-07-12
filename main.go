package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// type FileMetadata struct {
// 	SourcePath string
// 	TargetPath string
// }

func GetFilePaths(source, target string) ([]FileMetadata, error) {
	var files []FileMetadata
	var lastWeek = time.Now().Add(-24 * 7 * time.Hour)
	var err = filepath.Walk(source, func(path string, info os.FileInfo, err1 error) error {
		if err1 != nil {
			return err1
		}
		var relativePath, err2 = filepath.Rel(source, path)
		if err2 != nil {
			return err2
		}
		if !info.IsDir() && info.ModTime().After(lastWeek) {
			var metadata = FileMetadata{
				SourcePath: path,
				TargetPath: fmt.Sprintf("%s/%s", target, relativePath),
			}
			files = append(files, metadata)
		}
		return nil
	})
	if err != nil {
		return nil, err
	} else {
		return files, nil
	}
}

func SaveMetadataToFile(files []FileMetadata) error {
	var json, err1 = json.MarshalIndent(files, "", "\t")
	if err1 != nil {
		return err1
	}
	var err2 = ioutil.WriteFile("files.json", json, 0644)
	if err2 != nil {
		return err2
	}
	return nil
}

func main() {
	var cwd, err1 = os.Getwd()
	if err1 != nil {
		log.Fatal(err1)
	}
	var files, err2 = GetFilePaths("/home/applegamer22/Documents/scr-web/storage", cwd)
	if err2 != nil {
		log.Fatal(err2)
	}
	var err3 = SaveMetadataToFile(files)
	if err3 != nil {
		log.Fatal(nil)
	}
}
