package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

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

