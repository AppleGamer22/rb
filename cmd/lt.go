package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/AppleGamer22/recursive-backup/internal/manager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	ltCmd.PersistentFlags().StringVarP(&timeString, "time", "t", "", "reference time with format: 20060102T150405")
	viper.BindPFlag("reference_time", ltCmd.Flags().Lookup("time"))
	rootCmd.AddCommand(ltCmd)
}

var ltCmd = &cobra.Command{
	Use:   "lt",
	Short: "list recent modifications in source",
	Long:  "list source recursively for recent modifications in directories and files",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("arguments mismatch, expecting 1 argument")
		}
		cfg.Source = args[0]

		if timeString == "" {
			return errors.New("time string cannot be empty")
		}
		assertedTime, err := parseTime(timeString)
		if err != nil {
			return fmt.Errorf("failed to parse time flag value: %v", err)
		}
		cfg.ReferenceTime = assertedTime

		return nil
	},
	PreRunE: lsCmd.PreRunE,
	RunE:    ltRunCmd,
}

func ltRunCmd(cmd *cobra.Command, args []string) error {
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
		SourceRootDir: cfg.Source,
	}

	service := manager.NewService(in)
	if err = service.ListSources(dirs, files, errs, cfg.ReferenceTime); err != nil {
		return err
	}

	operationLogLine = "list end"
	if err = writeOpLog(operationLogLine); err != nil {
		return err
	}

	skeletonFormatString := "Run the following from the command line in order to create directories on the target directory:\n" +
		"\t%s skeleton -d \"%s\" -p \"%s\" \"%s\" \"[target-dir-path]\"\n"
	fmt.Printf(skeletonFormatString, os.Args[0], listDirsPath, rootDirPath, cfg.Source)
	sliceFormatString := "\nThen, run the following from the command line in order to divide the workload into smaller chunks:\n" +
		"\t%s slice -f \"%s\" -p \"%s\" -s [positive--integer-batch-size]\n"
	fmt.Printf(sliceFormatString, os.Args[0], listFilesPath, rootDirPath)
	return nil
}
