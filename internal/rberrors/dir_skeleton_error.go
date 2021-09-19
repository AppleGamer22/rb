package rberrors

import "fmt"

type DirSkeletonError struct {
	MissedDirPaths []string
}

func (d DirSkeletonError) Error() string {
	return fmt.Sprintf("missed directories: %v", d.MissedDirPaths)
}
