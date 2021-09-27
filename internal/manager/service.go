package manager

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/AppleGamer22/recursive-backup/internal/rberrors"
	"github.com/AppleGamer22/recursive-backup/internal/tasks"
	"github.com/AppleGamer22/recursive-backup/internal/workers"
)

type API interface {
	ListSources() error
	CreateTargetDirSkeleton() error
	RequestFilesCopy(filesList io.Reader)
	HandleFilesCopyResponse()
}

type service struct {
	// orientation
	SourceRootDir string
	TargetRootDir string

	// source listing dependencies
	ListingDirPathsWriter  io.Writer
	ListingFilePathsWriter io.Writer
	ListingErrorsLogWriter io.Writer

	// directories skeleton dependencies
	DirPathsReader io.Reader

	// files backup dependencies
	FilePathsReader     io.Reader
	TasksPipeline       chan tasks.BackupFile
	ResponsesChannel    chan tasks.BackupFileResponse
	FileBackupWorkers   []workers.FileBackupWorker
	FileBackupLogWriter io.Writer
}

type ServiceInitInput struct {
	SourceRootDir          string
	TargetRootDir          string
	ListingDirPathsWriter  io.Writer
	ListingFilePathsWriter io.Writer
	ListingErrorsLogWriter io.Writer
	FilePathsReader        io.Reader
	FileCopyPipelineLength int
	FileBackupLogWriter    io.Writer
}

//func (i ServiceInitInput) Validate() error {
//	return validation.ValidateStruct(&i,
//		validation.Field(&i.SourceRootDir, validation.Required, validation.By(val.CheckDirReadable)),
//		validation.Field(&i.TargetRootDir, validation.Required, validation.By(val.CheckDirReadable)),
//		validation.Field(&i.ListingDirPathsWriter, validation.Required),
//		validation.Field(&i.ListingFilePathsWriter, validation.Required),
//		validation.Field(&i.ListingErrorsLogWriter, validation.Required),
//		validation.Field(&i.FilePathsReader, validation.Required),
//		validation.Field(&i.FileCopyPipelineLength, validation.Required),
//		validation.Field(&i.FileBackupLogWriter, validation.Required),
//	)
//}

func NewService(in ServiceInitInput) API {
	tasksPipelineChannel := make(chan tasks.BackupFile, in.FileCopyPipelineLength)
	fileBackupResponsesChannel := make(chan tasks.BackupFileResponse, in.FileCopyPipelineLength)
	fileBackupWorkers := initFileWorkers(in.SourceRootDir, in.TargetRootDir, in.FileCopyPipelineLength, tasksPipelineChannel)

	return &service{
		SourceRootDir:          in.SourceRootDir,
		TargetRootDir:          in.TargetRootDir,
		ListingDirPathsWriter:  in.ListingDirPathsWriter,
		ListingFilePathsWriter: in.ListingFilePathsWriter,
		ListingErrorsLogWriter: in.ListingErrorsLogWriter,
		FilePathsReader:        nil,
		TasksPipeline:          tasksPipelineChannel,
		ResponsesChannel:       fileBackupResponsesChannel,
		FileBackupWorkers:      fileBackupWorkers,
		FileBackupLogWriter:    in.FileBackupLogWriter,
	}
}

func initFileWorkers(srcRootDir, targetRootDir string, numFileCopyWorkers int, bc chan tasks.BackupFile) []workers.FileBackupWorker {
	fileWorkers := make([]workers.FileBackupWorker, numFileCopyWorkers)
	for i := range fileWorkers {
		fileWorkers[i] = workers.NewFileBackupWorker(srcRootDir, targetRootDir, bc)
	}
	return fileWorkers
}

func (m *service) ListSources() error {
	newSourceListerInput := &tasks.NewSrcListerInput{
		SrcRootDir:   m.SourceRootDir,
		DirsWriter:   m.ListingDirPathsWriter,
		FilesWriter:  m.ListingFilePathsWriter,
		ErrorsWriter: m.ListingErrorsLogWriter,
	}

	sourceLister, err := tasks.NewSourceLister(newSourceListerInput)
	if err != nil {
		return err
	}

	return sourceLister.Do()
}

func (m *service) CreateTargetDirSkeleton() error {
	bufferedErrorsWriter := bufio.NewWriter(m.ListingErrorsLogWriter)
	task := tasks.NewBackupDirSkeleton(m.SourceRootDir, m.DirPathsReader, m.TargetRootDir)
	errs := task.Do()
	for _, err := range errs {
		switch err.(type) {
		case rberrors.DirSkeletonError:
			for _, missedPath := range err.(rberrors.DirSkeletonError).MissedDirPaths {
				msg := fmt.Sprintf("%s missed-path: %s\n", "dir-skeleton-error", missedPath)
				bufferedErrorsWriter.WriteString(msg)
			}
		default:
			msg := fmt.Sprintf("%s general-error: %s\n", "dir-skeleton-error", err.Error())
			bufferedErrorsWriter.WriteString(msg)
		}
	}
	if len(errs) > 0 {
		return errors.New("CreateTargetDirSkeleton completed with errors")
	}
	return nil
}

func (m *service) RequestFilesCopy(filesList io.Reader) {
	const replaceCount = 1
	buf := bufio.NewScanner(filesList)
	for buf.Scan() {
		srcPath := buf.Text()
		task := tasks.BackupFile{
			CreationRequestTime: time.Now(),
			SourcePath:          srcPath,
			TargetPath:          strings.Replace(srcPath, m.SourceRootDir, m.TargetRootDir, replaceCount),
			ResponseChannel:     m.ResponsesChannel,
		}
		task.Do()
	}
}

func (m *service) HandleFilesCopyResponse() {
	buf := bufio.NewWriter(m.FileBackupLogWriter)
	headerLine := fmt.Sprintf("status,duration [milli-sec],target,source\n") //todo: relocate
	buf.WriteString(headerLine)
	for resp := range m.ResponsesChannel {
		buf.WriteString(fileCopyResponseString(resp))
	}
}

func fileCopyResponseString(r tasks.BackupFileResponse) string {
	duration := r.CompletionTime.Sub(r.CreationRequestTime).Milliseconds()
	return fmt.Sprintf("%t,%d,%s,%s\n", r.CompletionStatus, duration, r.TargetPath, r.SourcePath)
}
