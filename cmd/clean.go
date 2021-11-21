package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	cleanCmd.Flags().StringVarP(&batchesDirPath, "batches-dir-path", "b", "", "mandatory flag: copy batches directory path")
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean [source-dir-path] [target-dir-path]",
	Short: "clean leftover files",
	Long:  "clean leftover files from batches to-do directory",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(batchesDirPath) == 0 {
			return errors.New("rootDirPath must be specified")
		}
		batchesToDoDirPath = filepath.Join(batchesDirPath, sliceBatchesToDoDirName)
		batchesDoneDirPath = filepath.Join(batchesDirPath, sliceBatchesDoneDirName)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS != "windows" {
			return nil
		}
		doneDirWalkFunc := func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				toBeRemovedPath := filepath.Join(batchesToDoDirPath, d.Name())
				if err := os.Remove(toBeRemovedPath); err != nil {
					_ = writeOpLog(fmt.Sprintf("%s (%+v)", toBeRemovedPath, err))
				}
			}
			return nil
		}

		_ = filepath.WalkDir(batchesDoneDirPath, doneDirWalkFunc)

		return nil
	},
}
