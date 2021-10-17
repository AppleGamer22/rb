package cmd

import (
	"fmt"
	"github.com/AppleGamer22/recursive-backup/internal/manager"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

var listDirPath = filepath.Join(rootDirPath, listDirName)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list all source elements",
	Long:  "list source recursively for all directories and files",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("arguments mismatch, expecting 1 argument")
		}
		cfg.Src = args[0]

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := initCmd.RunE(cmd, args); err != nil {
			return err
		}
		if err :=  os.Chdir(listDirPath); err != nil {
			return err
		}
		return nil
	},
	RunE: listRunCommand,
}

func listRunCommand(cmd *cobra.Command, args []string) error {
	fmt.Printf("src: %v\n", cfg.Src)

	operationLogLine := fmt.Sprintf("list start %s", time.Now().String())
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

	operationLogLine = fmt.Sprintf("list end %s", time.Now().String())
	if err := writeOpLog(operationLogLine); err != nil {
		return err
	}

	return nil
}

func createFilesForList() (dirs, files, errs *os.File, err error) {
	dirs, err = os.Create(fmt.Sprintf(listedDirsFileNamePattern, time.Now().Format(timeDateFormat)))
	if err != nil {
		return nil, nil, nil, err
	}

	files, err = os.Create(fmt.Sprintf(listedFilesFileNamePattern, time.Now().Format(timeDateFormat)))
	if err != nil {
		return nil, nil, nil, err
	}

	errs, err = os.Create(fmt.Sprintf(listErrorsFileNamePattern, time.Now().Format(timeDateFormat)))
	if err != nil {
		return nil, nil, nil, err
	}

	return dirs, files, errs, nil
}

