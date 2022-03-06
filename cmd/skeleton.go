package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/AppleGamer22/recursive-backup/internal/manager"
	"github.com/AppleGamer22/recursive-backup/internal/rberrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var skeletonWorkDir string

func init() {
	skeletonCmd.Flags().StringVarP(&cfg.ProjectDir, "project", "p", "", "mandatory flag: project root path")
	skeletonCmd.Flags().StringVarP(&cfg.DirsListPath, "dirs-list-file-path", "d", "", "mandatory flag: directories list file path")
	skeletonCmd.Flags().StringVarP(&cfg.DirValidationMode, "dir-validation-mode", "v", rberrors.Report, "validation mode for directories short list (none, report, block)")

	viper.BindPFlag("ProjectDir", skeletonCmd.Flags().Lookup("project"))
	viper.BindPFlag("DirsListPath", skeletonCmd.Flags().Lookup("dirs-list-file-path"))
	viper.BindPFlag("DirValidationMode", skeletonCmd.Flags().Lookup("dir-validation-mode"))
	rootCmd.AddCommand(skeletonCmd)

}

var skeletonCmd = &cobra.Command{
	Use:   "skeleton [source-dir-path] [target-dir-path]",
	Short: "create target directory skeleton",
	Long:  "create directory skeleton in target",
	Args: func(cmd *cobra.Command, args []string) error {
		if cfg.DirValidationMode != rberrors.None && cfg.DirValidationMode != rberrors.Report && cfg.DirValidationMode != rberrors.Block {
			return fmt.Errorf("--on-missing-dir flag can be on of none, report or block, got %s", cfg.DirValidationMode)
		}

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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		skeletonWorkDir = filepath.Join(cfg.ProjectDir, dirSkeletonDirName)
		return nil
	},
	RunE: skeletonRunCommand,
}

func skeletonRunCommand(cmd *cobra.Command, args []string) error {
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
		SourceRootDir: cfg.Source,
		TargetRootDir: cfg.Target,
	}
	service := manager.NewService(in)
	var reader io.Reader
	if reader, err = service.CreateTargetDirSkeleton(inDirsListFile, errorsFile, cfg.DirValidationMode); err != nil {
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
	inDirsList, err = os.Open(cfg.DirsListPath)
	if err != nil {
		return nil, nil, nil, err
	}

	skeletonDirsFileName := fmt.Sprintf(skeletonDirsFileNamePattern, time.Now().Format(timeDateFormat))
	outDirsList, err = os.Create(skeletonDirsFileName)
	if err != nil {
		return nil, nil, nil, err
	}
	fmt.Println(filepath.Join(skeletonWorkDir, skeletonDirsFileName))

	skeletonDirErrorsFile := fmt.Sprintf(skeletonErrorsFileNamePattern, time.Now().Format(timeDateFormat))
	errs, err = os.Create(skeletonDirErrorsFile)
	if err != nil {
		return nil, nil, nil, err
	}
	fmt.Println(filepath.Join(skeletonWorkDir, skeletonDirErrorsFile))

	return inDirsList, outDirsList, errs, nil
}
