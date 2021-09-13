package tasks

import (
	"fmt"
	"io"
	"os"
	"time"
)

type BackupFile struct {
	CreationRequestTime time.Time
	SourcePath          string
	TargetPath          string
	ResponseChannel     chan BackupFileResponse
}

type BackupFileResponse struct {
	CreationRequestTime time.Time
	CompletionTime      time.Time
	SourcePath          string
	TargetPath          string
	CompletionStatus    bool
}

func (b *BackupFile) Do() {
	_, err := copyFile(b.SourcePath, b.TargetPath)
	completionTime := time.Now()
	defer func() {
		b.ResponseChannel <- BackupFileResponse{
			CreationRequestTime: b.CreationRequestTime,
			CompletionTime:      completionTime,
			SourcePath:          b.SourcePath,
			TargetPath:          b.TargetPath,
			CompletionStatus: func() bool {
				if err == nil {
					return true
				}
				return false
			}(),
		}
	}()
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)

	return nBytes, err
}
