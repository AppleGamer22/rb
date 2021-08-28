package workers

import "github.com/AppleGamer22/recursive-backup/internal/tasks"

type FileBackupWorker interface {
	Do() error
}

type fileBackupWorker struct {
	pipeline chan tasks.BackupFile
}

func NewFileBackupWorker(p chan tasks.BackupFile) FileBackupWorker {
	return fileBackupWorker{
		pipeline: p,
	}
}

func (f fileBackupWorker) Do() error{

	return nil
}
