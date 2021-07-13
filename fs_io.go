package main

import (
	"encoding/json"
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

func RemoveMetadataFromFile(path string) error {
	var data, err1 = ioutil.ReadFile(path)
	if err1 != nil {
		return err1
	}
	var metadata []FileMetadata
	var err2 = json.Unmarshal(data, &metadata)
	if err2 != nil {
		return err2
	}
	println(len(metadata))
	if len(metadata) > 0 {
		metadata = metadata[1:]
	}
	println(len(metadata))
	var err3 = SaveMetadataToFile(metadata, path)
	if err3 != nil {
		return err3
	}
	return nil
}