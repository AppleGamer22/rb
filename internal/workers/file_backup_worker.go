package workers

import (
	"github.com/AppleGamer22/recursive-backup/internal/tasks"
)

type FileBackupWorker interface {
	Handle()
}

type fileBackupWorker struct {
	ID             uint
	Pipeline       chan tasks.GeneralRequest
	SourceRootPath string
	TargetRootPath string
	QuitFunc       UpdateOnQuitFunc
}

type UpdateOnQuitFunc func()

func NewFileBackupWorker(id uint, srcRootPath, targetRootPath string, p chan tasks.GeneralRequest, quitFunc UpdateOnQuitFunc) {
	worker := &fileBackupWorker{
		ID:             id,
		Pipeline:       p,
		SourceRootPath: srcRootPath,
		TargetRootPath: targetRootPath,
		QuitFunc:       quitFunc,
	}
	go worker.Handle()
}

func (f *fileBackupWorker) Handle() {
	for task := range f.Pipeline {
		switch assertedRequest := task.(type) {
		case tasks.BackupFileRequest:
			assertedRequest.WorkerID = f.ID
			response := assertedRequest.Do()
			assertedRequest.ResponseChannel <- response
		case tasks.QuitRequest:
			f.QuitFunc()
			return
		default:
			return
		}

	}
}
