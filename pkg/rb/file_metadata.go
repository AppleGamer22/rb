package rb

import "time"

type Status string
const (
	Success Status = "success"
	Fail Status = "fail"
	Unchanged Status = "unchanged"
)

type FileMetadata struct {
	SourcePath string `json:"source"`
	TargetPath string `json:"target"`
	CompletionStatus bool `json:"status"`
	ChangedStatus bool `json:"changes"`
}

type Log struct {
	LastBackupTime time.Time `json:"lastBackupTime"`
	FilesCopied int `json:"filesCopied"`
	Files []FileMetadata `json:"files"`
}