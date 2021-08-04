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

type RecursiveBackupper struct {
	SourceRoot string
	TargetRoot string
	StartTime  time.Time
	PreviousExecutionTime *time.Time
	RecoveryMode bool
}

func NewRecursiveBackupper(sourceRoot string, targetRoot string, previousExecutionTime *time.Time, recover bool) RecursiveBackupper {
	return RecursiveBackupper{
		SourceRoot: sourceRoot,
		TargetRoot: targetRoot,
		StartTime: time.Now(),
		PreviousExecutionTime: previousExecutionTime,
		RecoveryMode: recover,
	}
}

// Copies all files changed after provided time from source directory to target directory.
func (rber RecursiveBackupper) BackupFilesSinceDate() (executionLogPath string, err error) {
	executionLogPath = fmt.Sprintf("rb_%d-%d-%d_%d.%d.%d.csv", rber.StartTime.Day(), rber.StartTime.Month(), rber.StartTime.Year(), rber.StartTime.Hour(), rber.StartTime.Minute(), rber.StartTime.Second())
	file, err := os.Create(executionLogPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	var count = 0
	err = filepath.Walk(rber.SourceRoot, func(sourcePath string, info os.FileInfo, err error) error {
		utils.WaitForDirectory(rber.SourceRoot)
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
			foundOnTarget, err := utils.DoesTargetFileExist(sourcePath, rber.SourceRoot, rber.TargetRoot)
			if err != nil {
				return err
			}
			isInitialBackup := rber.PreviousExecutionTime == nil
			isRecent := rber.PreviousExecutionTime != nil && info.ModTime().After(*rber.PreviousExecutionTime)
			isRecoverable := rber.RecoveryMode && !foundOnTarget
			if rber.RecoveryMode && foundOnTarget {
				fmt.Printf("Skipping %s\n", sourcePath)
			}
			if isInitialBackup || isRecent || isRecoverable {
				targetFilePath, copyTime, err := rber.backupFile(sourcePath, count)
				if err != nil {
					writer.WriteString(fmt.Sprintf("ERROR: %s, %s\n", sourcePath, strings.ReplaceAll(err.Error(), "\n", "")))
				}
				writer.WriteString(fmt.Sprintf("%s,%s,%s\n", sourcePath, targetFilePath, copyTime))
				count++
			}
		} else if info.Mode().IsDir() {
			targetPath, err := utils.Source2TargetPath(sourcePath, rber.SourceRoot, rber.TargetRoot)
			if err != nil {
				return err
			}
			err = os.MkdirAll(targetPath, 0755)
			if err != nil {
				return err
			}
		}
		return nil
	})
	writer.Flush()
	return executionLogPath, err
}