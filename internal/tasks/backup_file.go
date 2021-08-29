package tasks

import "time"

type BackupFile struct {
	CreationTime time.Time
	SourcePath   string
	TargetPath   string
}
