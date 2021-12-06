package manager

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/AppleGamer22/recursive-backup/internal/rberrors"
	"github.com/AppleGamer22/recursive-backup/internal/tasks"
	val "github.com/AppleGamer22/recursive-backup/internal/validationhelpers"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type API interface {
	ListSources(dirsWriter, filesWriter, errorsWriter io.Writer) error
	ListSourcesReferenceTime(dirsWriter, filesWriter, errorsWriter io.Writer) error
	CreateTargetDirSkeleton(dirsReader io.Reader, errorsWriter io.Writer) (io.Reader, error)
	RequestFilesCopy(filesList io.Reader, batchID uint, requestChan chan tasks.GeneralRequest, responseChan chan tasks.BackupFileResponse)
	HandleFilesCopyResponse(logWriter io.Writer, responseChan chan tasks.BackupFileResponse)
	WaitForAllResponses()
}

type service struct {
	SourceRootDir         string
	TargetRootDir         string
	RecoveryReferenceTime time.Time
}

type ServiceInitInput struct {
	SourceRootDir         string
	TargetRootDir         string
	RecoveryReferenceTime time.Time
}

var wgRequestResponseCorelator sync.WaitGroup

func (i ServiceInitInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.SourceRootDir, validation.Required, validation.By(val.CheckDirReadable)),
		validation.Field(&i.TargetRootDir, validation.Required, validation.By(val.CheckDirReadable)),
	)
}

func NewService(in ServiceInitInput) API {
	return &service{
		SourceRootDir: in.SourceRootDir,
		TargetRootDir: in.TargetRootDir,
	}
}

func (m *service) ListSources(dirsWriter, filesWriter, errorsWriter io.Writer) error {
	newSourceListerInput := &tasks.NewSrcListerInput{
		SrcRootDir:   m.SourceRootDir,
		DirsWriter:   dirsWriter,
		FilesWriter:  filesWriter,
		ErrorsWriter: errorsWriter,
	}

	sourceLister, err := tasks.NewSourceLister(newSourceListerInput)
	if err != nil {
		return err
	}

	return sourceLister.Do(false)
}

func (m *service) ListSourcesReferenceTime(dirsWriter, filesWriter, errorsWriter io.Writer) error {
	newSourceListerInput := &tasks.NewSrcListerInput{
		SrcRootDir:            m.SourceRootDir,
		RecoveryReferenceTime: m.RecoveryReferenceTime,
		DirsWriter:            dirsWriter,
		FilesWriter:           filesWriter,
		ErrorsWriter:          errorsWriter,
	}

	sourceLister, err := tasks.NewSourceLister(newSourceListerInput)
	if err != nil {
		return err
	}

	return sourceLister.Do(true)
}

func (m *service) CreateTargetDirSkeleton(srcDirsReader io.Reader, errorsWriter io.Writer) (io.Reader, error) {
	bufferedErrorsWriter := bufio.NewWriter(errorsWriter)
	task := tasks.NewBackupDirSkeleton(srcDirsReader, m.SourceRootDir, m.TargetRootDir)
	createdDirsReader, errs := task.Do()
	for _, err := range errs {
		switch err.(type) {
		case rberrors.DirSkeletonError:
			for _, missedPath := range err.(rberrors.DirSkeletonError).MissedDirPaths {
				msg := fmt.Sprintf("%s missed-path: %s\n", "dir-skeleton-error", missedPath)
				_, _ = bufferedErrorsWriter.WriteString(msg)
			}
		default:
			msg := fmt.Sprintf("%s general-error: %s\n", "dir-skeleton-error", err.Error())
			_, _ = bufferedErrorsWriter.WriteString(msg)
		}
	}
	_ = bufferedErrorsWriter.Flush()
	if len(errs) > 0 {
		return createdDirsReader, errors.New("CreateTargetDirSkeleton completed with errors")
	}
	return createdDirsReader, nil
}

func (m *service) RequestFilesCopy(filesList io.Reader, batchID uint, requestChan chan tasks.GeneralRequest, responseChan chan tasks.BackupFileResponse) {
	scanner := bufio.NewScanner(filesList)
	var fileID uint = 0
	for scanner.Scan() {
		srcFullPath := scanner.Text()
		filePath := strings.TrimPrefix(srcFullPath, m.SourceRootDir)
		targetFullPath := filepath.Join(m.TargetRootDir, filePath)
		copyFileTask := tasks.BackupFileRequest{
			FileID:              fileID,
			BatchID:             batchID,
			CreationRequestTime: time.Now(),
			SourcePath:          srcFullPath,
			TargetPath:          targetFullPath,
			ResponseChannel:     responseChan,
		}
		requestChan <- copyFileTask
		wgRequestResponseCorelator.Add(1)
		fileID++
	}
}

func (m *service) HandleFilesCopyResponse(logWriter io.Writer, responseChan chan tasks.BackupFileResponse) {
	buf := bufio.NewWriter(logWriter)
	headerLine := "status,duration [milli-sec],target,source,error_message\n"
	_, _ = buf.WriteString(headerLine)
	for resp := range responseChan {
		_, _ = buf.WriteString(fileCopyResponseString(resp))
		_ = buf.Flush()
		wgRequestResponseCorelator.Done()
	}
}

func fileCopyResponseString(r tasks.BackupFileResponse) string {
	duration := r.CompletionTime.Sub(r.CreationRequestTime).Milliseconds()
	return fmt.Sprintf("%t,%d,%s,%s,%s\n", r.CompletionStatus, duration, r.TargetPath, r.SourcePath, r.ErrorMessage)
}

func (m *service) WaitForAllResponses() {
	wgRequestResponseCorelator.Wait()
}
