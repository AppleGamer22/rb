package manager

import (
	"io"

	"github.com/AppleGamer22/recursive-backup/internal/tasks"
	"github.com/AppleGamer22/recursive-backup/internal/workers"
)

type Manager interface {
	CreateTargetDirSkeleton(dirList io.Reader) error
	ProcessFilesCopyRequest(filesList io.Reader) error
}

type manager struct {
	sourceRootDir           string
	targetRootDir           string
	tasksPipeline           chan tasks.BackupFile
	directorySkeletonWorker workers.DirectorySkeletonWorker
	fileBackupWorkers       []workers.FileBackupWorker
}

func NewManager(srcRootDir string, targetRootDir string, pipelineLength int) Manager {
	pipelineChannel := make(chan tasks.BackupFile, pipelineLength)

	fileWorkers := make([]workers.FileBackupWorker, pipelineLength)
	for i := range fileWorkers {
		fileWorkers[i] = workers.NewFileBackupWorker(pipelineChannel)
	}

	return manager{
		sourceRootDir:           srcRootDir,
		targetRootDir:           targetRootDir,
		tasksPipeline:           pipelineChannel,
		directorySkeletonWorker: workers.NewDirectorySkeletonWorker(),
		fileBackupWorkers:       fileWorkers,
	}
}

func (m manager) CreateTargetDirSkeleton(dirList io.Reader) error {

	return nil
}

func (m manager) ProcessFilesCopyRequest(filesList io.Reader) error {

	return nil
}
