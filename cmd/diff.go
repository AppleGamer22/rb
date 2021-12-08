package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var timeString string

func init() {
	diffCmd.Flags().UintVarP(&cfg.NumWorkers, "workers", "w", 2, "number of concurrent file-copy workers")
	diffCmd.Flags().StringVarP(&timeString, "time", "t", "", "reference time with format: 2006-01-02T15:04:05")
	rootCmd.AddCommand(diffCmd)
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "differential backup",
	Long:  "with differential backup only new or modified files and directories are being backed-up",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("arguments mismatch, expecting 2 arguments")
		}
		cfg.Src = args[0]
		cfg.Target = args[1]

		assertedTime, err := time.Parse(timeFormat, timeString)
		if err != nil {
			return fmt.Errorf("failed to parse time flag value")
		}
		if !assertedTime.Before(time.Now()) {
			return fmt.Errorf("reference time flag value is in the future")
		}
		cfg.RecoveryReferenceTime = assertedTime

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// list
		listDirPath = filepath.Join(rootDirPath, listDirName)
		if err := ltCmd.RunE(cmd, args); err != nil {
			return err
		}

		// skeleton
		skeletonWorkDir = filepath.Join(rootDirPath, dirSkeletonDirName)
		dirsListFilePath = listDirsPath
		if err := skeletonCmd.RunE(cmd, args); err != nil {
			return err
		}

		// slice
		filesListFilePath = listFilesPath
		if err := sliceCmd.PreRunE(cmd, args); err != nil {
			return err
		}
		if err := sliceCmd.RunE(cmd, args); err != nil {
			return err
		}

		// cp
		batchesDirPath = batchesSourceDirPath
		if err := cpCmd.RunE(cmd, args); err != nil {
			return err
		}

		return nil
	},
}
