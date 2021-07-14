package main

import "time"

type FileMetadata struct {
	SourcePath string `json:"source"`
	TargetPath string `json:"target"`
	Done bool `json:"done"`
}

type Log struct {
	LastBackupTime time.Time `json:"lastBackupTime"`
	FilesCopied int `json:"filesCopied"`
	Files []FileMetadata `json:"files"`
}