package workers

import (
	"github.com/AppleGamer22/recursive-backup/internal/tasks"
)

type Copy interface {
	Handle()
}

type copyWorker struct {
	ID             uint
	Pipeline       chan tasks.GeneralRequest
	SourceRootPath string
	TargetRootPath string
	QuitFunc       UpdateOnQuitFunc
}

type UpdateOnQuitFunc func()

func NewCopyWorker(id uint, srcRootPath, targetRootPath string, p chan tasks.GeneralRequest, quitFunc UpdateOnQuitFunc) {
	worker := &copyWorker{
		ID:             id,
		Pipeline:       p,
		SourceRootPath: srcRootPath,
		TargetRootPath: targetRootPath,
		QuitFunc:       quitFunc,
	}
	go worker.Handle()
}

func (f *copyWorker) Handle() {
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
