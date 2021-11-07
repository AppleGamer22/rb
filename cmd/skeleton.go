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

var skeletonWorkDir string
var dirsListFilePath string

func init() {
	skeletonCmd.Flags().StringVarP(&rootDirPath, "project", "p", "", "project root path")
	skeletonCmd.Flags().StringVarP(&dirsListFilePath, "dirs-list-file-path", "d", "", "directories list file path")
	rootCmd.AddCommand(skeletonCmd)

}

var skeletonCmd = &cobra.Command{
	Use:   "skeleton",
	Short: "create target directory skeleton",
	Long:  "create directory skeleton in target",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("arguments mismatch, expecting 2 arguments")
		}
		cfg.Src = args[0]
		cfg.Target = args[1]

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		skeletonWorkDir = filepath.Join(rootDirPath, dirSkeletonDirName)
		return nil
	},
	RunE: skeletonRunCommand,
}

func skeletonRunCommand(cmd *cobra.Command, args []string) error {
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
