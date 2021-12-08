package rberrors

import "fmt"

const (
	None   = "none"
	Report = "report"
	Block  = "block"
)

type DirSkeletonError struct {
	MissedDirPaths []string
}

func (d DirSkeletonError) Error() string {
	return fmt.Sprintf("missed directories: %v", d.MissedDirPaths)
}
