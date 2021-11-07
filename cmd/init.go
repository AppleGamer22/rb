package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/spf13/cobra"
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
	slicesWorkDirName             = "slices"
	copyLogDirPattern             = "copy" + string(filepath.Separator) + "copy_logs_%s"
	copyBatchLogFileNamePattern   = "copy_batch_%d.log"
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
			_ = writeOpLog("init error")
			return err
		}
		_ = writeOpLog("init successful finish")
		return nil
	},
}

func setup() error {
	err := validateWorkDir(false)
	if err != nil {
		return err
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
	_ = writeOpLog("init start")

	subDirs := []string{listDirName, dirSkeletonDirName, slicesWorkDirName}
	for _, subDir := range subDirs {
		if err = os.Mkdir(subDir, defaultPerm); err != nil {
			return fmt.Errorf("failed to create directory %s", subDir)
		}
	}
	return nil
}

func validateWorkDir(enforceProjectDir bool) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	re := regexp.MustCompile(parentDirNameRegexp)
	if re == nil {
		return errors.New("invalid regexp")
	}
	if !re.Match([]byte(wd)) {
		if enforceProjectDir {
			return errors.New("this command need to be executed from within a project directory")
		}
	}
	rootDirPath = wd
	return nil
}
