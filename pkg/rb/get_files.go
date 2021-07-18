package rb

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AppleGamer22/recursive-backup/pkg/utils"
)

// Copies all files changed after provided time from source directory to target directory.
func BackupFilesSinceDate(source, target string, date *time.Time) (string, time.Time, error) {
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
		utils.WaitForDirectory(source)
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
		if info.Mode().IsRegular() {
			if (date != nil && info.ModTime().After(*date)) || date == nil {
				sourceFilePath, targetFilePath, copyTime, err := BackupFile(sourcePath, source, target, count, now)
				if err != nil {
					writer.WriteString(fmt.Sprintf("ERROR: %s, %s\n", sourcePath, strings.ReplaceAll(err.Error(), "\n", "")))
				}
				writer.WriteString(fmt.Sprintf("%s,%s,%s\n", sourceFilePath, targetFilePath, copyTime))
				fmt.Println("done")
				count++
			}
		} else if info.Mode().IsDir() {
			_, err := os.Stat(sourcePath)
			if err != nil {
				targetPath, err := utils.Source2TargetPath(sourcePath, source, target)
				if err != nil {
					return err
				}
				err = os.MkdirAll(targetPath, 0755)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	writer.Flush()
	return logsPath, now, err
}