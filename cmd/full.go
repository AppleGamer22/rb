package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	// slice dependency
	fullCmd.Flags().UintVarP(&batchSize, "batch-size", "s", defaultBatchSize, "maximum number of files in a batch")

	// cp dependency
	fullCmd.Flags().UintVarP(&copyQueueLen, "copy-queue-len", "q", 200, "copy queue length")

	rootCmd.AddCommand(fullCmd)
}

var fullCmd = &cobra.Command{
	Use:   "full [source-dir-path] [target-dir-path]",
	Short: "full backup",
	Long:  "with full backup all files and folders are copied from src to target",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("arguments mismatch, expecting 2 arguments: [source-dir-path] [target-dir-path]")
		}
		cfg.Src = args[0]
		cfg.Target = args[1]
		return nil
	},
	PreRunE: initCmd.RunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		// list
		listDirPath = filepath.Join(rootDirPath, listDirName)
		if err := lsCmd.RunE(cmd, args); err != nil {
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
		if err := cpCmd.PreRunE(cmd, args); err != nil {
			return err
		}
		if err := cpCmd.RunE(cmd, args); err != nil {
			return err
		}

		return nil
	},
}
