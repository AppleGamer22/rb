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
	responseChannel         chan tasks.BackupFileResponse
	directorySkeletonWorker workers.DirectorySkeletonWorker
	fileBackupWorkers       []workers.FileBackupWorker
}

func NewManager(srcRootDir string, targetRootDir string, pipelineLength int) Manager {
	pipelineChannel := make(chan tasks.BackupFile, pipelineLength)
	filesCopyResponseChannel := make(chan tasks.BackupFileResponse, pipelineLength)

	fileWorkers := make([]workers.FileBackupWorker, pipelineLength)
	for i := range fileWorkers {
		fileWorkers[i] = workers.NewFileBackupWorker(srcRootDir, targetRootDir, pipelineChannel)
	}

	return manager{
		sourceRootDir:           srcRootDir,
		targetRootDir:           targetRootDir,
		tasksPipeline:           pipelineChannel,
		responseChannel:         filesCopyResponseChannel,
		directorySkeletonWorker: workers.NewDirectorySkeletonWorker(srcRootDir, targetRootDir),
		fileBackupWorkers:       fileWorkers,
	}
}

func (m manager) CreateTargetDirSkeleton(dirList io.Reader) error {

	return nil
}

func (m manager) ProcessFilesCopyRequest(filesList io.Reader) error {

	return nil
}
