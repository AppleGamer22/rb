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

func GetFilePaths(source, target string) (string, time.Time, error) {
	var now = time.Now()
	var executionLogPath = fmt.Sprintf("rb_%d-%d-%d_%d.%d.%d.csv", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
	file, err := os.Create(executionLogPath)
	if err != nil {
		return "", now, err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	var count = 0
	err = filepath.Walk(source, func(sourcePath string, info os.FileInfo, err error) error {
		WaitForDirectory(source)
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
			sourceFilePath, targetFilePath, modTime, err := BackupFile(sourcePath, executionLogPath, source, target, count, now)
			if err != nil {
				writer.WriteString(fmt.Sprintf("ERROR: %s, %s\n", sourcePath, strings.ReplaceAll(err.Error(), "\n", "")))
			}
			writer.WriteString(fmt.Sprintf("%s,%s,%s\n", sourceFilePath, targetFilePath, modTime))
			fmt.Println("done")
		}
		return nil
	})
	writer.Flush()
	return executionLogPath, now, err
}

func GetFilePathsSinceDate(source, target string, date time.Time) (string, time.Time, error) {
	var now = time.Now()
	var logsPath = fmt.Sprintf("rb_%d-%d-%d_%d.%d.%d.csv", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
	file, err := os.Create(logsPath)
	if err != nil {
		return "", now, err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	var count = 0
	err = filepath.Walk(source, func(sourcePath string, info os.FileInfo, err error) error {
		WaitForDirectory(source)
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
			sourceFilePath, targetFilePath, modTime, err := BackupFile(sourcePath, logsPath, source, target, count, now)
			if err != nil {
				writer.WriteString(fmt.Sprintf("ERROR: %s, %s\n", sourcePath, strings.ReplaceAll(err.Error(), "\n", "")))
			}
			writer.WriteString(fmt.Sprintf("%s,%s,%s\n", sourceFilePath, targetFilePath, modTime))
			fmt.Println("done")
			count++
		}
		return nil
	})
	writer.Flush()
	return logsPath, now, err
}