package rb

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetFilePaths(source, target string) (string, int, time.Time, error) {
	var now = time.Now()
	var sourcesLogPath = fmt.Sprintf("rb_source_%d-%d-%d_%d.%d.%d.csv", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
	file, err := os.Create(sourcesLogPath)
	if err != nil {
		return "", 0, now, err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	var count = 0
	err = filepath.Walk(source, func(sourcePath string, info os.FileInfo, err error) error {
		if err != nil {
			writer.WriteString(fmt.Sprintf("ERROR: %s, %s\n", sourcePath, strings.ReplaceAll(err.Error(), "\n", "")))
			switch err.(type) {
			case *fs.PathError:
				return nil
			default:
				fmt.Println(err)
				fmt.Printf("%T\n", err)
				return err
			}
		}
		if !info.IsDir() {
			_, err := writer.WriteString(sourcePath + "\n")
			if err != nil {
				return err
			}
			count++
			fmt.Println(sourcePath)
		}
		return nil
	})
	writer.Flush()
	return sourcesLogPath, count, now, err
}

func GetFilePathsSinceDate(source, target string, date time.Time) (string, int, time.Time, error) {
	var now = time.Now()
	var sourcesLogPath = fmt.Sprintf("rb_source_%d-%d-%d_%d.%d.%d.csv", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
	file, err := os.Create(sourcesLogPath)
	if err != nil {
		return "", 0, now, err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	var count = 0
	err = filepath.Walk(source, func(sourcePath string, info os.FileInfo, err error) error {
		if err != nil {
			writer.WriteString(fmt.Sprintf("ERROR: %s, %s\n", sourcePath, strings.ReplaceAll(err.Error(), "\n", "")))
			switch err.(type) {
			case *fs.PathError:
				return nil
			default:
				fmt.Println(err)
				fmt.Printf("%T\n", err)
				return err
			}
		}
		if !info.IsDir() && info.ModTime().After(date) {
			_, err := writer.WriteString(sourcePath + "\n")
			if err != nil {
				return err
			}
			count++
			fmt.Println(sourcePath)
		}
		return nil
	})
	writer.Flush()
	return sourcesLogPath, count, now, err
}
