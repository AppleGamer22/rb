package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	diffCmd.Flags().UintVarP(&cfg.NumWorkers, "copy-queue-len", "q", 200, "copy queue length")
	diffCmd.Flags().UintVarP(&cfg.BatchSize, "batch-size", "s", defaultBatchSize, "maximum number of files in a batch")
	diffCmd.Flags().StringVarP(&cfg.ReferenceTimeString, "time", "t", "", "reference time with format: 20060102T150405")

	viper.BindPFlag("num_workers", diffCmd.Flags().Lookup("copy-queue-len"))
	viper.BindPFlag("batch_size", diffCmd.Flags().Lookup("batch-size"))
	viper.BindPFlag("reference_time", diffCmd.Flags().Lookup("time"))
	rootCmd.AddCommand(diffCmd)
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "differential backup",
	Long:  "with differential backup only new or modified files and directories are being backed-up",
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

		if cfg.ReferenceTimeString == "" {
			return errors.New("time string cannot be empty")
		}
		_, err := parseTime(cfg.ReferenceTimeString)
		if err != nil {
			return fmt.Errorf("failed to parse time flag value: %v", err)
		}

		return nil
	},
	PreRunE: fullCmd.PreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		// list
		listDirPath = filepath.Join(cfg.ProjectDir, listDirName)
		if err := ltCmd.RunE(cmd, args); err != nil {
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
