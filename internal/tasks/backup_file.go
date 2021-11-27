package tasks

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type GeneralRequest interface{}

type QuitRequest struct{}

type BackupFileRequest struct {
	WorkerID            uint
	FileID              uint
	BatchID             uint
	CreationRequestTime time.Time
	SourcePath          string
	TargetPath          string
	ResponseChannel     chan BackupFileResponse
}

type BackupFileResponse struct {
	WorkerID            uint
	FileID              uint
	BatchID             uint
	CreationRequestTime time.Time
	CompletionTime      time.Time
	SourcePath          string
	TargetPath          string
	CompletionStatus    bool
	ErrorMessage        string
}

func (b *BackupFileRequest) Do() BackupFileResponse {
	fmt.Printf(">>[w%d][b%d][f%d]>> cp %s -> %s\n", b.WorkerID, b.BatchID, b.FileID, b.SourcePath, b.TargetPath)
	_, err := copyFile(b.SourcePath, b.TargetPath)
	switch err.(type) {
	case *fs.PathError:
		dirPath := filepath.Dir(b.TargetPath)
		err = os.MkdirAll(dirPath, 0755)
		if err == nil {
			_, err = copyFile(b.SourcePath, b.TargetPath)
		}
	}

	response := BackupFileResponse{
		WorkerID:            b.WorkerID,
		BatchID:             b.BatchID,
		FileID:              b.FileID,
		CreationRequestTime: b.CreationRequestTime,
		CompletionTime:      time.Now(),
		SourcePath:          b.SourcePath,
		TargetPath:          b.TargetPath,
		CompletionStatus:    err == nil,
		ErrorMessage: func() string {
			var val = "success"
			if err != nil {
				val = err.Error()
			}
			return val
		}(),
	}

	fmt.Printf("<<[w%d][b%d][f%d][%t]<< cp %s -> %s\n", b.WorkerID,
		b.BatchID, b.FileID, response.CompletionStatus, b.SourcePath, b.TargetPath)

	return response
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
	defer func() {
		_ = source.Close()
	}()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = destination.Close()
	}()

	nBytes, err := io.Copy(destination, source)

	return nBytes, err
}
