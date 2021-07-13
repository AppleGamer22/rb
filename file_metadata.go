package main

type FileMetadata struct {
	SourcePath string `json:"source"`
	TargetPath string `json:"target"`
	Done bool `json:"done"`
}