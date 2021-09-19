package tasks

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	val "github.com/AppleGamer22/recursive-backup/internal/validationhelpers"
	"github.com/go-ozzo/ozzo-validation/v4"
)

type SourceLister interface {
	Do() error
}

type sourceLister struct {
	SrcRootDir   string
	DirsWriter   *bufio.Writer
	FilesWriter  *bufio.Writer
	ErrorsWriter *bufio.Writer
}

type NewSrcListerInput struct {
	SrcRootDir   string
	DirsWriter   io.Writer
	FilesWriter  io.Writer
	ErrorsWriter io.Writer
}

func (i *NewSrcListerInput) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.SrcRootDir, validation.Required, validation.By(val.CheckDirReadable)),
		validation.Field(&i.DirsWriter, validation.Required, validation.NotNil),
		validation.Field(&i.FilesWriter, validation.Required, validation.NotNil),
		validation.Field(&i.ErrorsWriter, validation.Required, validation.NotNil),
	)
}

func NewSourceLister(input *NewSrcListerInput) (SourceLister, error) {
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

func (s *sourceLister) Do() error {
	return filepath.WalkDir(s.SrcRootDir, s.walkDirFunc)
}

func (s *sourceLister) walkDirFunc(path string, d fs.DirEntry, err error) error {
	switch {
	case err != nil:
		_, err = s.ErrorsWriter.WriteString(fmt.Sprintf("%s, %v\n", path, err))
		if err != nil {
			return err
		}
		if d.IsDir() {
			return fs.SkipDir
		}
	case d.IsDir():
		_, err = s.DirsWriter.WriteString(fmt.Sprintf("%s\n", path))
		if err != nil {
			return err
		}
	case dirEntryError(d) != nil:
		_, err = s.ErrorsWriter.WriteString(fmt.Sprintf("%s, %v\n", path, dirEntryError(d)))
		if err != nil {
			return err
		}
	case isRegular(d):
		_, err = s.FilesWriter.WriteString(fmt.Sprintf("%s\n", path))
		if err != nil {
			return err
		}
	default:
		err = errors.New("unexpected_element")
		_, err = s.ErrorsWriter.WriteString(fmt.Sprintf("%s, %v\n", path, err))
		if err != nil {
			return err
		}
	}
	return nil
}

func isRegular(d fs.DirEntry) bool {
	fileInfo, _ := d.Info()
	return fileInfo.Mode().IsRegular()
}

func dirEntryError(d fs.DirEntry) error {
	_, err := d.Info()
	return err
}
