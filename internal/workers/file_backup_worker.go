package workers

import (
	"github.com/AppleGamer22/recursive-backup/internal/tasks"
)

type FileBackupWorker interface {
	Handle()
}

type fileBackupWorker struct {
	Pipeline       chan tasks.BackupFile
	SourceRootPath string
	TargetRootPath string
}

func NewFileBackupWorker(srcRootPath, targetRootPath string, p chan tasks.BackupFile) FileBackupWorker {
	worker := &fileBackupWorker{
		Pipeline:       p,
		SourceRootPath: srcRootPath,
		TargetRootPath: targetRootPath,
	}
	go worker.Handle()
	return worker
}

func (f *fileBackupWorker) Handle() {
	for task := range f.Pipeline {
		task.Do()
	}
}
