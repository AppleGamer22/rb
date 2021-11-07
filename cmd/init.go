package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	parentDirNameRegexp           = ".*" + string(filepath.Separator) + "rb_[[:digit:]]{8}T[[:digit:]]{6}$"
	timeDateFormat                = "20060102T150405"
	parentDirNamePattern          = "rb_%s"
	listDirName                   = "list"
	dirSkeletonDirName            = "dirs"
	sliceBatchesDirNamePattern    = "slice" + string(filepath.Separator) + "batches_%s"
	sliceBatchesErrorDirPattern   = "slice" + string(filepath.Separator) + "errors_%s"
	listedDirsFileNamePattern     = "list_dirs_%s.log"
	listedFilesFileNamePattern    = "list_files_%s.log"
	listErrorsFileNamePattern     = "list_errors_%s.log"
	skeletonDirsFileNamePattern   = "skeleton_dirs_%s.log"
	skeletonErrorsFileNamePattern = "skeleton_errors_%s.log"
	sliceBatchFileNamePattern     = "batch_%d.log"
	sliceErrorsFileNamePattern    = "slice_errors_%s.log"
	fileTasksDirName              = "file_batches"
	copyLogDirPattern             = "file_copy_%s" + string(filepath.Separator) + "copy_log"
	copyErrorLogDirPattern        = "file_copy_%s" + string(filepath.Separator) + "error_log"
	copyBatchLogFileNamePattern   = "copy_batch_%d.log"
	copyErrorLogFileName          = "copy_error.log"
	errorsDirName                 = "errors"
	operationLogFileName          = "oplog.log"
	defaultPerm                   = 0755
)

var rootDirPath string

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init rb project",
	Long:  "init initialised a new backup project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := setup(); err != nil {
			return err
		}
		return nil
	},
}

func setup() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	re := regexp.MustCompile(parentDirNameRegexp)
	if re == nil {
		return errors.New("invalid regexp")
	}
	if re.Match([]byte(wd)) {
		rootDirPath = wd
		return nil
	}

	rootDirName := fmt.Sprintf(parentDirNamePattern, time.Now().Format(timeDateFormat))
	if err = os.Mkdir(rootDirName, defaultPerm); err != nil {
		return errors.New("failed to create patent directory")
	}
	if rootDirPath, err = filepath.Abs(rootDirName); err != nil {
		return errors.New("failed to create absolute path for work dir")
	}
	if err = os.Chdir(rootDirName); err != nil {
		return fmt.Errorf("failed to change work dir to %s", rootDirName)
	}

	operationLogLine := "init"
	if err = writeOpLog(operationLogLine); err != nil {
		return err
	}

	subDirs := []string{listDirName, dirSkeletonDirName, fileTasksDirName, errorsDirName}
	for _, subDir := range subDirs {
		if err = os.Mkdir(subDir, defaultPerm); err != nil {
			return fmt.Errorf("failed to create directory %s", subDir)
		}
	}
	return nil
}
