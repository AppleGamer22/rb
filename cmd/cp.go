package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/AppleGamer22/recursive-backup/internal/manager"
	"github.com/spf13/cobra"
)

var cpWorkDir string
var fileCopyFilePath string

func init() {
	cpCmd.Flags().StringVarP(&rootDirPath, "project", "p", "", "project root path")
	cpCmd.Flags().StringVarP(&dirsListFilePath, "file-copy-list", "f", "", "file path to a file with the list of files to copy")
	rootCmd.AddCommand(cpCmd)

}

var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "copy files",
	Long:  "copy files from source to target dir",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("arguments mismatch, expecting 2 arguments")
		}
		cfg.Src = args[0]
		cfg.Target = args[1]

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cpWorkDir = filepath.Join(rootDirPath, dirSkeletonDirName)
		return nil
	},
	RunE: cpRunCommand,
}

func cpRunCommand(cmd *cobra.Command, args []string) error {
	fmt.Printf("src: %v\n", cfg.Src)

	operationLogLine := "directory skeleton build start"
	if err := writeOpLog(operationLogLine); err != nil {
		return err
	}

	if err := os.Chdir(skeletonWorkDir); err != nil {
		return err
	}

	inDirsListFile, outDirsListFile, errorsFile, err := setupForDirSkeleton()
	if err != nil {
		return err
	}
	defer func() {
		_ = inDirsListFile.Close()
		_ = outDirsListFile.Close()
		_ = errorsFile.Close()
	}()

	in := manager.ServiceInitInput{
		SourceRootDir: cfg.Src,
		TargetRootDir: cfg.Target,
	}
	service := manager.NewService(in)
	var reader io.Reader
	if reader, err = service.CreateTargetDirSkeleton(inDirsListFile, errorsFile); err != nil {
		return err
	}
	if _, err = io.Copy(outDirsListFile, reader); err != nil {
		return err
	}

	operationLogLine = "directory skeleton build end"
	if err = writeOpLog(operationLogLine); err != nil {
		return err
	}

	return nil
}

func setupForDirSkeleton() (inDirsList, outDirsList, errs *os.File, err error) {
	inDirsList, err = os.Open(dirsListFilePath)
	if err != nil {
		return nil, nil, nil, err
	}

	outDirsList, err = os.Create(fmt.Sprintf(skeletonDirsFileNamePattern, time.Now().Format(timeDateFormat)))
	if err != nil {
		return nil, nil, nil, err
	}

	errs, err = os.Create(fmt.Sprintf(skeletonErrorsFileNamePattern, time.Now().Format(timeDateFormat)))
	if err != nil {
		return nil, nil, nil, err
	}

	return inDirsList, outDirsList, errs, nil
}
