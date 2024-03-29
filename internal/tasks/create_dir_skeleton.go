package tasks

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/AppleGamer22/recursive-backup/internal/rberrors"
)

type BackupDirSkeleton interface {
	Do() (io.Reader, []error)
}

func NewBackupDirSkeleton(srcDirReader io.Reader, srcRootPath string, targetRootPath string, validationMode string) BackupDirSkeleton {
	return &backupDirSkeleton{
		SrcRootPath:          srcRootPath,
		SrcDirectoriesReader: srcDirReader,
		ValidationMode:       validationMode,
		TargetRootPath:       targetRootPath,
	}
}

type backupDirSkeleton struct {
	SrcRootPath          string
	SrcDirectoriesReader io.Reader
	ValidationMode       string
	TargetRootPath       string
}

func (b *backupDirSkeleton) Do() (io.Reader, []error) {
	var errs []error
	dirs, err := b.extractLongPaths()
	if err != nil && b.ValidationMode == "block" {
		errs = append(errs, err)
		return nil, errs
	}

	builder := strings.Builder{}
	for _, srcDirPath := range dirs {
		trimmedSrcDirPath := strings.TrimPrefix(srcDirPath, b.SrcRootPath)
		expression := fmt.Sprintf("\\%c{2,}", filepath.Separator)
		exp := regexp.MustCompile(expression)
		targetDirPath := fmt.Sprintf("%s%c%s", b.TargetRootPath, filepath.Separator, trimmedSrcDirPath)
		targetDirPath = exp.ReplaceAllString(targetDirPath, string(filepath.Separator))

		err = os.MkdirAll(targetDirPath, 0755)
		if err != nil {
			errs = append(errs, err)
		} else {
			builder.WriteString(fmt.Sprintf("%s\n", targetDirPath))
		}
	}

	paths := strings.Split(builder.String(), "\n")
	sort.Strings(paths)
	out := strings.Join(paths, "\n")
	out = strings.TrimPrefix(out, "\n")

	return strings.NewReader(out), errs
}

func (b *backupDirSkeleton) extractLongPaths() (shortList []string, err error) {
	dirs, err := b.getSortedSrcDirPaths()
	if err != nil {
		return []string{}, err
	}

	// var shortList []string
	var lastDir string
	for _, dir := range dirs {
		if !strings.HasPrefix(lastDir, dir) {
			shortList = append(shortList, dir)
			lastDir = dir
		}
	}

	var missed []string
	for _, dirPath := range dirs {
		fileInfo, err := os.Stat(dirPath)
		if err != nil || !fileInfo.IsDir() {
			switch b.ValidationMode {
			case "report":
			case "block":
				missed = append(missed, dirPath)
			}
		}
	}
	if len(missed) > 0 {
		err = rberrors.DirSkeletonError{MissedDirPaths: missed}
	}

	return shortList, err
}

func (b *backupDirSkeleton) getSortedSrcDirPaths() ([]string, error) {
	var dirs []string

	scan := bufio.NewScanner(b.SrcDirectoriesReader)
	for scan.Scan() {
		text := scan.Text()
		dirs = append(dirs, text)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(dirs)))

	return dirs, nil
}
