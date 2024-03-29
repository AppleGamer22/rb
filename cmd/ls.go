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
	if err = service.ListSources(dirs, files, errs, nil); err != nil {
		return err
	}

	operationLogLine = "list end"
	if err = writeOpLog(operationLogLine); err != nil {
		return err
	}

	skeletonFormatString := "Run the following from the command line in order to create directories on the target directory:\n" +
		"\t%s skeleton -d \"%s\" -p \"%s\" \"%s\" \"[target-dir-path]\"\n"
	fmt.Printf(skeletonFormatString, os.Args[0], listDirsPath, rootDirPath, cfg.Src)
	sliceFormatString := "\nThen, run the following from the command line in order to divide the workload into smaller chunks:\n" +
		"\t%s slice -f \"%s\" -p \"%s\" -s [positive--integer-batch-size]\n"
	fmt.Printf(sliceFormatString, os.Args[0], listFilesPath, rootDirPath)
	return nil
}

func createFilesForList() (dirs, files, errs *os.File, err error) {
	now := time.Now()
	listDirsName := fmt.Sprintf(listedDirsFileNamePattern, now.Format(timeDateFormat))
	listDirsPath = filepath.Join(listDirPath, listDirsName)
	dirs, err = os.Create(listDirsPath)
	if err != nil {
		return nil, nil, nil, err
	}
	fmt.Println(listDirsPath)

	listFilesName := fmt.Sprintf(listedFilesFileNamePattern, now.Format(timeDateFormat))
	listFilesPath = filepath.Join(listDirPath, listFilesName)
	files, err = os.Create(listFilesPath)
	if err != nil {
		return nil, nil, nil, err
	}
	fmt.Println(listFilesPath)

	errorsFileName := fmt.Sprintf(listErrorsFileNamePattern, now.Format(timeDateFormat))
	listErrorsPath := filepath.Join(listDirPath, errorsFileName)
	errs, err = os.Create(listErrorsPath)
	if err != nil {
		return nil, nil, nil, err
	}
	fmt.Println(listErrorsPath)

	return dirs, files, errs, nil
}
