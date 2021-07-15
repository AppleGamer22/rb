package rb

import "time"


type FileMetadata struct {
	SourcePath string `json:"source"`
	TargetPath string `json:"target"`
	CompletionStatus bool `json:"status"`
	ChangedStatus bool `json:"changes"`
}

type SourceLog struct {
	BackupStartTime time.Time `json:"startTime"`
	SourceFiles []string `json:"files"`
}

type Log struct {
	LastBackupTime time.Time `json:"lastBackupTime"`
	FilesCopied int `json:"filesCopied"`
	Files []FileMetadata `json:"files"`
}