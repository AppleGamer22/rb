package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AppleGamer22/recursive-backup/pkg/rb"
	"github.com/AppleGamer22/recursive-backup/pkg/utils"
)

// Detects previous execution log if it is provided and exists and starts the back-up process.
func PrepareData(source, target string, logsFlag *string, recoveryFlag *bool) (string, error) {
	fmt.Printf("Source directory: %s\n", source)
	fmt.Printf("Target directory: %s\n", target)
	if (*recoveryFlag) {
		fmt.Printf("Copying files not found on %s\n", target)
	}
	previousExecutionLogPath := *logsFlag
	if previousExecutionLogPath != "" {
		_, err := os.Stat(previousExecutionLogPath)
		if err != nil {
			fmt.Println("could not get logs file data from ", previousExecutionLogPath, "\n Error: ", err)
			os.Exit(1)
		}
		previousExecutionTime, _ := utils.GetLastBackupExecutionTime(previousExecutionLogPath)
		rber := rb.NewRecursiveBackupper(source, target, &previousExecutionTime, *recoveryFlag)
		path, err := rber.BackupFilesSinceDate()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		return path, nil
	} else {
		rber := rb.NewRecursiveBackupper(source, target, nil, *recoveryFlag)
		path, err := rber.BackupFilesSinceDate()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		return path, nil
	}
}

// Accepts arguments/flags and starts the program.
func main() {
	var logsFlag = flag.String("logs", "", "--logs \"<logs JSON file path>\"")
	var recoveryFlag = flag.Bool("recovery", false, "--recovery")
	flag.Parse()
	if len(flag.Args()) == 0 {
		showHelp()
		return
	}
	source, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		fmt.Println("source path is invalid")
		os.Exit(1)
	}

	target, err := filepath.Abs(flag.Arg(1))
	if err != nil {
		fmt.Println("target path is not valid")
		os.Exit(1)
	}
	newLogsFilePath, err := PrepareData(source, target, logsFlag, recoveryFlag)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("The Backup log is saved at: %s", newLogsFilePath)
}

// prints usage guide to console.
func showHelp() {
	fmt.Println("Usage:")
	fmt.Println("For full backup:")
	fmt.Println("\trb \"<source path>\" \"<target path>\"")
	fmt.Println("For partial backup:")
	fmt.Println("\trb --logs \"<logs JSON file path>\" \"<source path>\" \"<target path>\"")
	fmt.Println("For usage guide:")
	fmt.Println("\trb")
}
