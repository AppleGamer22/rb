package tasks

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/AppleGamer22/recursive-backup/internal/rberrors"
)

type BackupDirSkeleton struct {
	SrcRootPath          string
	SrcDirectoriesReader io.Reader
	TargetRootPath       string
}

func (b *BackupDirSkeleton) Do() []error {
	var errs []error

	dirs, err := b.extractLongPaths()
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	for _, srcDirPath := range dirs {
		targetDirPath := b.TargetRootPath + strings.TrimPrefix(srcDirPath, b.SrcRootPath)
		err = os.MkdirAll(targetDirPath, 0755)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (b *BackupDirSkeleton) extractLongPaths() ([]string, error) {
	dirs, err := b.getSortedSrcDirPaths()
	if err != nil {
		return nil, err
	}

	var shortList []string
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
			missed = append(missed, dirPath)
		}
	}
	if len(missed) > 0 {
		err = rberrors.DirSkeletonError{MissedDirPaths: missed}
	}

	return shortList, err
}

func (b *BackupDirSkeleton) getSortedSrcDirPaths() ([]string, error) {
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
