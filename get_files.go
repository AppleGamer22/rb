package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func GetFilePaths(source, target string) ([]FileMetadata, error) {
	var files []FileMetadata
	// var lastWeek = time.Now().Add(-24 * 7 * time.Hour)
	var err = filepath.Walk(source, func(path string, info os.FileInfo, err1 error) error {
		if err1 != nil {
			return err1
		}
		var relativePath, err2 = filepath.Rel(source, path)
		if err2 != nil {
			return err2
		}
		if !info.IsDir() {
			var metadata = FileMetadata{
				SourcePath: path,
				TargetPath: fmt.Sprintf("%s/%s", target, relativePath),
				Done: false,
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

func GetFilePathsSinceDate(source, target string, date time.Time) ([]FileMetadata, error) {
	var files []FileMetadata
	// var lastWeek = time.Now().Add(-24 * 7 * time.Hour)
	var err = filepath.Walk(source, func(path string, info os.FileInfo, err1 error) error {
		if err1 != nil {
			return err1
		}
		var relativePath, err2 = filepath.Rel(source, path)
		if err2 != nil {
			return err2
		}
		if !info.IsDir() && info.ModTime().After(date) {
			var metadata = FileMetadata{
				SourcePath: path,
				TargetPath: fmt.Sprintf("%s/%s", target, relativePath),
				Done: false,
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