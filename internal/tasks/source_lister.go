package tasks

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type SourceLister interface {
	Do() error
}

type sourceLister struct {
	SrcRootDir   string
	DirsWriter   bufio.Writer
	FilesWriter  bufio.Writer
	ErrorsWriter bufio.Writer
}

func (s *sourceLister) Validate() error {
	checkSrcDir := func(srcDir interface{}) error {
		assertedSrcDir, ok := srcDir.(string)
		if !ok {
			return errors.New("must be a string")
		}
		fileInfo, err := os.Stat(assertedSrcDir)
		if err != nil {
			return errors.New("must be accessible")
		}
		if !fileInfo.IsDir() {
			return errors.New("must be a directory path")
		}
		_, err = os.ReadDir(assertedSrcDir)
		if err != nil {
			return errors.New("must be readable")
		}
		return nil
	}

	return validation.ValidateStruct(s,
		validation.Field(s.SrcRootDir, validation.Required, validation.By(checkSrcDir)),
		validation.Field(s.FilesWriter, validation.Required),
		validation.Field(s.DirsWriter, validation.Required),
		validation.Field(s.ErrorsWriter, validation.Required),
	)
}

func NewSourceLister(srcRootDir string, dirsWriter, filesWriter, errorsWriter bufio.Writer) (SourceLister, error) {
	srcLister := &sourceLister{
		SrcRootDir:   srcRootDir,
		DirsWriter:   dirsWriter,
		FilesWriter:  filesWriter,
		ErrorsWriter: errorsWriter,
	}

	err := srcLister.Validate()

	return srcLister, err
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
