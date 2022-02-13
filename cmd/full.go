package cmd

import (
	"errors"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// slice dependency
	fullCmd.Flags().UintVarP(&cfg.BatchSize, "batch-size", "s", defaultBatchSize, "maximum number of files in a batch")

	// cp dependency
	fullCmd.Flags().UintVarP(&cfg.NumWorkers, "copy-queue-len", "q", 200, "copy queue length")

	viper.BindPFlag("batch_size", fullCmd.Flags().Lookup("batch-size"))
	viper.BindPFlag("num_workers", fullCmd.Flags().Lookup("copy-queue-len"))
	rootCmd.AddCommand(fullCmd)
}

var fullCmd = &cobra.Command{
	Use:   "full [source-dir-path] [target-dir-path]",
	Short: "full backup",
	Long:  "with full backup all files and folders are copied from src to target",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			if len(cfg.Source) == 0 && len(cfg.Target) == 0 {
				return errors.New("arguments mismatch, expecting 2 arguments: [source-dir-path] [target-dir-path]")
			} else if (len(cfg.Source) > 0 || len(cfg.Target) > 0) && !(len(cfg.Source) > 0 && len(cfg.Target) > 0) {
				return errors.New("arguments mismatch, expecting 2 arguments: [source-dir-path] [target-dir-path]")
			}
		} else {
			cfg.Source = args[0]
			cfg.Target = args[1]
		}
		return nil
	},
	PreRunE: initCmd.RunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		// list
		listDirPath = filepath.Join(cfg.ProjectDir, listDirName)
		if err := lsCmd.RunE(cmd, args); err != nil {
			return err
		}

		// skeleton
		skeletonWorkDir = filepath.Join(cfg.ProjectDir, dirSkeletonDirName)
		cfg.DirsListPath = listDirsPath
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
		cfg.BatchesDirPath = batchesSourceDirPath
		if err := cpCmd.PreRunE(cmd, args); err != nil {
			return err
		}
		if err := cpCmd.RunE(cmd, args); err != nil {
			return err
		}

		return nil
	},
}
