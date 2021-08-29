package tasks

import "io"

type CreateDirSkeleton struct{
	DirectoryList io.Reader
}
