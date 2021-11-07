package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AppleGamer22/recursive-backup/internal/manager"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

var listDirPath string
var listFilesPath string
var listDirsPath string

var lsCmd = &cobra.Command{
	Use:   "ls [source-dir-path]",
	Short: "list all source elements",
	Long:  "list source recursively for all directories and files",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("arguments mismatch, expecting 1 argument: [source-dir-path]")
		}
		cfg.Src = args[0]

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateWorkDir(true); err != nil {
			if err = initCmd.RunE(cmd, args); err != nil {
				return err
			}
		}

		listDirPath = filepath.Join(rootDirPath, listDirName)

		if err := os.Chdir(listDirPath); err != nil {
			return err
		}
		return nil
	},
	RunE: listRunCommand,
}

func listRunCommand(cmd *cobra.Command, args []string) error {
	operationLogLine := "list start"
	if err := writeOpLog(operationLogLine); err != nil {
		return err
	}

	if err := os.Chdir(listDirPath); err != nil {
		return err
	}

	dirs, files, errs, err := createFilesForList()
	if err != nil {
		return err
	}
	defer func() {
		_ = dirs.Close()
		_ = files.Close()
		_ = errs.Close()
	}()

	in := manager.ServiceInitInput{
		SourceRootDir: cfg.Src,
	}
	service := manager.NewService(in)
	if err = service.ListSources(dirs, files, errs); err != nil {
		return err
	}

	operationLogLine = "list end"
	if err = writeOpLog(operationLogLine); err != nil {
		return err
	}

	return nil
}

func createFilesForList() (dirs, files, errs *os.File, err error) {
	now := time.Now()
	listDirsName := fmt.Sprintf(listedDirsFileNamePattern, now.Format(timeDateFormat))
	dirs, err = os.Create(listDirsName)
	if err != nil {
		return nil, nil, nil, err
	}
	listDirsPath = filepath.Join(listDirPath, listDirsName)

	listFilesName := fmt.Sprintf(listedFilesFileNamePattern, now.Format(timeDateFormat))
	files, err = os.Create(listFilesName)
	if err != nil {
		return nil, nil, nil, err
	}
	listFilesPath = filepath.Join(listDirPath, listFilesName)

	errs, err = os.Create(fmt.Sprintf(listErrorsFileNamePattern, now.Format(timeDateFormat)))
	if err != nil {
		return nil, nil, nil, err
	}

	return dirs, files, errs, nil
}
