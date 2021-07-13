package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func Copy(file FileMetadata) error {
	fileStat, err := os.Stat(file.SourcePath)
	if err != nil {
		return err
	}
	if !fileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", file.SourcePath)
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

func SaveMetadataToFile(files []FileMetadata, path string) error {
	var json, err1 = json.MarshalIndent(files, "", "\t")
	if err1 != nil {
		return err1
	}
	var err2 = ioutil.WriteFile(path, json, 0644)
	if err2 != nil {
		return err2
	}
	return nil
}

func MarkAsDone(path string, i int) error {
	var data, err1 = ioutil.ReadFile(path)
	if err1 != nil {
		return err1
	}
	var metadata []FileMetadata
	var err2 = json.Unmarshal(data, &metadata)
	if err2 != nil {
		return err2
	}
	if len(metadata) > 0 && 0 <= i && i < len(metadata) {
		metadata[i].Done = true
	} else {
		return errors.New("index is out of scope")
	}
	var err3 = SaveMetadataToFile(metadata, path)
	if err3 != nil {
		return err3
	}
	return nil
}
