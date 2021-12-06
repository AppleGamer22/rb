package tasks

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"time"

	val "github.com/AppleGamer22/recursive-backup/internal/validationhelpers"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SourceListerAPI interface {
	Do(withReferenceTime bool) error
}

type sourceLister struct {
	SrcRootDir            string
	RecoveryReferenceTime time.Time
	DirsWriter            *bufio.Writer
	FilesWriter           *bufio.Writer
	ErrorsWriter          *bufio.Writer
}

type NewSrcListerInput struct {
	SrcRootDir            string
	RecoveryReferenceTime time.Time
	DirsWriter            io.Writer
	FilesWriter           io.Writer
	ErrorsWriter          io.Writer
}

func (i *NewSrcListerInput) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.SrcRootDir, validation.Required, validation.By(val.CheckDirReadable)),
		validation.Field(&i.DirsWriter, validation.Required, validation.NotNil),
		validation.Field(&i.FilesWriter, validation.Required, validation.NotNil),
		validation.Field(&i.ErrorsWriter, validation.Required, validation.NotNil),
	)
}

func NewSourceLister(input *NewSrcListerInput) (SourceListerAPI, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	srcLister := &sourceLister{
		SrcRootDir:   input.SrcRootDir,
		DirsWriter:   bufio.NewWriter(input.DirsWriter),
		FilesWriter:  bufio.NewWriter(input.FilesWriter),
		ErrorsWriter: bufio.NewWriter(input.ErrorsWriter),
	}

	return srcLister, nil
}

func (s *sourceLister) Do(withReferenceTime bool) error {
	if withReferenceTime {
		if err := filepath.WalkDir(s.SrcRootDir, s.walkDirFuncReferenceTime); err != nil {
			return err
		}
	} else {
		if err := filepath.WalkDir(s.SrcRootDir, s.walkDirFunc); err != nil {
			return err
		}
	}
	_ = s.DirsWriter.Flush()
	_ = s.FilesWriter.Flush()
	_ = s.ErrorsWriter.Flush()

	return nil
}

func (s *sourceLister) walkDirFunc(path string, d fs.DirEntry, err error) error {
	switch {
	case err != nil:
		_, _ = s.ErrorsWriter.WriteString(fmt.Sprintf("%s, %v\n", path, err))
		// if d.IsDir() {
		// 	return fs.SkipDir
		// }
	case d.IsDir():
		_, _ = s.DirsWriter.WriteString(fmt.Sprintf("%s\n", path))
	case dirEntryError(d) != nil:
		_, _ = s.ErrorsWriter.WriteString(fmt.Sprintf("%s, %v\n", path, dirEntryError(d)))
		return fs.SkipDir
	case isRegular(d):
		_, _ = s.FilesWriter.WriteString(fmt.Sprintf("%s\n", path))
	default:
		msg := "unexpected_element"
		_, _ = s.ErrorsWriter.WriteString(fmt.Sprintf("path: %s, type: %v error_msg: %s\n", path, d.Type(), msg))
	}
	return nil
}

func (s *sourceLister) walkDirFuncReferenceTime(path string, d fs.DirEntry, err error) error {
	switch {
	case err != nil:
		_, _ = s.ErrorsWriter.WriteString(fmt.Sprintf("%s, %v\n", path, err))
		// if d.IsDir() {
		// 	return fs.SkipDir
		// }
	case d.IsDir() && isAfterReferenceTime(d, s.RecoveryReferenceTime):
		_, _ = s.DirsWriter.WriteString(fmt.Sprintf("%s\n", path))
	case dirEntryError(d) != nil:
		_, _ = s.ErrorsWriter.WriteString(fmt.Sprintf("%s, %v\n", path, dirEntryError(d)))
		return fs.SkipDir
	case isRegular(d) && isAfterReferenceTime(d, s.RecoveryReferenceTime):
		_, _ = s.FilesWriter.WriteString(fmt.Sprintf("%s\n", path))
	default:
		msg := "unexpected_element"
		_, _ = s.ErrorsWriter.WriteString(fmt.Sprintf("path: %s, type: %v error_msg: %s\n", path, d.Type(), msg))
	}
	return nil
}

func isRegular(d fs.DirEntry) bool {
	fileInfo, _ := d.Info()
	return fileInfo.Mode().IsRegular()
}

func isAfterReferenceTime(d fs.DirEntry, referenceTime time.Time) bool {
	fileInfo, _ := d.Info()
	return fileInfo.ModTime().After(referenceTime)
}

func dirEntryError(d fs.DirEntry) error {
	_, err := d.Info()
	return err
}
