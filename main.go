package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type FileMetadata struct {
	SourcePath string
	TargetPath string
}

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
	var sourceFlag = flag.String("source", "", "")
	var targetFlag = flag.String("target", "", "")
	flag.Parse()
	var source, err1 = filepath.Abs(string(*sourceFlag))
	if err1 != nil {
		log.Fatal("source path is invalid")
	}
	var target, err2 = filepath.Abs(string(*targetFlag))
	if err2 != nil {
		log.Fatal("target path is not valid")
	}
	var files, err3 = GetFilePaths(source, target)
	if err3 != nil {
		log.Fatal(err3)
	}
	var err4 = SaveMetadataToFile(files)
	if err4 != nil {
		log.Fatal(nil)
	}
}
