package rb

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func GetFilePaths(source, target string) (string, int, time.Time, error) {
	var now = time.Now()
	var sourcesLogPath = fmt.Sprintf("rb_source_%d-%d-%d_%d:%d:%d.txt", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
	file, err := os.Create(sourcesLogPath)
	if err != nil {
		return "", 0, now, err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	var count = 0
	err = filepath.Walk(source, func(sourcePath string, info os.FileInfo, err1 error) error {
		if err1 != nil {
			return err1
		}
		if !info.IsDir() {
			var _, err3 = writer.WriteString(sourcePath + "\n")
			if err3 != nil {
				return err3
			}
			writer.Flush()
			count++;
		}
		return nil
	})
	return sourcesLogPath, count, now, err
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
				CompletionStatus: false,
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