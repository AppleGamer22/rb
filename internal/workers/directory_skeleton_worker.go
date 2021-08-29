package workers

import "io"

type DirectorySkeletonWorker interface {
	Do(directoriesList io.Reader) error
}

type directorySkeletonWorker struct {}

func NewDirectorySkeletonWorker() DirectorySkeletonWorker{
	return directorySkeletonWorker{}
}

func (d directorySkeletonWorker) Do(directoriesList io.Reader) error{

	return nil
}
