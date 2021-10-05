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
	val "github.com/AppleGamer22/recursive-backup/internal/validationhelpers"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type API interface {
	ListSources(dirsWriter, filesWriter, errorsWriter io.Writer) error
	CreateTargetDirSkeleton(dirsReader io.Reader, errorsWriter io.Writer) (io.Reader, error)
	RequestFilesCopy(filesList io.Reader, responseChan chan tasks.BackupFileResponse)
	HandleFilesCopyResponse(logWriter io.Writer, responseChan chan tasks.BackupFileResponse)
}

type service struct {
	SourceRootDir string
	TargetRootDir string
}

type ServiceInitInput struct {
	SourceRootDir string
	TargetRootDir string
}

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

	return sourceLister.Do()
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
				bufferedErrorsWriter.WriteString(msg)
			}
		default:
			msg := fmt.Sprintf("%s general-error: %s\n", "dir-skeleton-error", err.Error())
			bufferedErrorsWriter.WriteString(msg)
		}
	}
	bufferedErrorsWriter.Flush()
	if len(errs) > 0 {
		return createdDirsReader, errors.New("CreateTargetDirSkeleton completed with errors")
	}
	return createdDirsReader, nil
}

func (m *service) RequestFilesCopy(filesList io.Reader, responseChan chan tasks.BackupFileResponse) {
	const replaceCount = 1
	buf := bufio.NewScanner(filesList)
	for buf.Scan() {
		srcPath := buf.Text()
		task := tasks.BackupFile{
			CreationRequestTime: time.Now(),
			SourcePath:          srcPath,
			TargetPath:          strings.Replace(srcPath, m.SourceRootDir, m.TargetRootDir, replaceCount),
			ResponseChannel:     responseChan,
		}
		task.Do()
	}
}

func (m *service) HandleFilesCopyResponse(logWriter io.Writer, responseChan chan tasks.BackupFileResponse) {
	buf := bufio.NewWriter(logWriter)
	headerLine := fmt.Sprintf("status,duration [milli-sec],target,source,error_message\n") //todo: relocate
	buf.WriteString(headerLine)
	for resp := range responseChan {
		buf.WriteString(fileCopyResponseString(resp))
		buf.Flush()
	}
}

func fileCopyResponseString(r tasks.BackupFileResponse) string {
	duration := r.CompletionTime.Sub(r.CreationRequestTime).Milliseconds()
	return fmt.Sprintf("%t,%d,%s,%s,%s\n", r.CompletionStatus, duration, r.TargetPath, r.SourcePath, r.ErrorMessage)
}
