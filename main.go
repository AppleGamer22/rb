package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/AppleGamer22/recursive-backup/pkg/rb"
)

func PrepareData(source, target string, logsFlag *string) (string, error) {
	logsPath := *logsFlag
	fmt.Println(logsPath) //DEBUG
	if logsPath != "" {
		stats, err := os.Stat(logsPath)
		if err != nil {
			fmt.Println("could not get logs file data from ", logsPath, "\n Error: ", err)
			os.Exit(1)
		}
		path, _, err := rb.GetFilePathsSinceDate(source, target, stats.ModTime())
		if err != nil {
			log.Fatal(err)
		}
		return path, nil
	} else {
		path, _, err := rb.GetFilePaths(source, target)
		if err != nil {
			log.Fatal(err)
		}
		return path, nil
	}
}

func main() {
	var logsFlag = flag.String("logs", "", "--logs \"<logs JSON file path>\"")
	flag.Parse()
	
	if len(flag.Args()) == 0 {
		ShowHelp()
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
	newLogsFilePath, err := PrepareData(source, target, logsFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The Backup log is saved at: %s", newLogsFilePath)
}

func ShowHelp() {
	fmt.Println("Usage:")
	fmt.Println("For full backup:")
	fmt.Println("\trb \"<source path>\" \"<target path>\"")
	fmt.Println("For partial backup:")
	fmt.Println("\trb --logs \"<logs JSON file path>\" \"<source path>\" \"<target path>\"")
	fmt.Println("For usage guide:")
	fmt.Println("\trb")
}
