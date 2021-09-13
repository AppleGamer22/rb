package workers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type DirectorySkeletonWorker interface {
	Do(directoriesList io.Reader) error
}

type directorySkeletonWorker struct {
	SourceRootPath string
	TargetRootPath string
}

func NewDirectorySkeletonWorker(srcRootPath, targetRootPath string) DirectorySkeletonWorker{
	return directorySkeletonWorker{
		SourceRootPath: srcRootPath,
		TargetRootPath: targetRootPath,
	}
}

func (d directorySkeletonWorker) Do(directoriesListReader io.Reader) error{
	scanner := bufio.NewScanner(directoriesListReader)
	scanner.Split(bufio.ScanLines)

	var failedSrcPaths []string
	for scanner.Scan() {
		srcPath := scanner.Text()
		targetPath := strings.Replace(srcPath, d.SourceRootPath, d.TargetRootPath, 1)
		if err := os.Mkdir(targetPath, 0755); err != nil {
			failedSrcPaths = append(failedSrcPaths, srcPath)
		}
	}

	if len(failedSrcPaths) == 0 {
		return nil
	} else {
		return fmt.Errorf("failed to create %q directories. Source direcory list: %v",
			len(failedSrcPaths), failedSrcPaths)
	}
}
