package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/AppleGamer22/recursive-backup/pkg/rb"
)

func main() {
	var logsFlag = flag.String("logs", "", "")
	flag.Parse()
	var source, err1 = filepath.Abs(flag.Arg(0))
	if err1 != nil {
		log.Fatal("source path is invalid")
	}
	var target, err2 = filepath.Abs(flag.Arg(1))
	if err2 != nil {
		log.Fatal("target path is not valid")
	}
	logsPath := string(*logsFlag)

	if len(flag.Args()) == 0 {
		ShowHelp()
		return
	}
	var now = time.Now()
	var fileCount int
	var newLogsPath = fmt.Sprintf("rb_%d-%d-%d_%d:%d:%d.csv", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
	var files []rb.FileMetadata
	if logsPath != "" {
		logs, err := rb.GetLogFromFile(logsPath)
		if err != nil {
			log.Fatal(err)
		}
		files, err = rb.GetFilePathsSinceDate(source, target, logs.LastBackupTime)
		if err != nil {
			log.Fatal(err)
		}
		var err4 = rb.SaveMetadataToFile(files, newLogsPath, 0, false, logs.LastBackupTime)
		if err4 != nil {
			log.Fatal(err4)
		}
	} else {
		var path, count, _, err2 = rb.GetFilePaths(source, target)
		if err2 != nil {
			log.Fatal(err2)
		}
		newLogsPath = path
		fileCount = count
	}
	var err = rb.Backup(newLogsPath, source, target, fileCount, now)
	if err != nil {
		log.Fatal(err)
	}
}

func ShowHelp() {
	fmt.Println("Usage:")
	fmt.Println("For full backup:")
	fmt.Println("\trb \"<source path>\" \"<target path>\"")
	fmt.Println("For partial backup:")
	fmt.Println("\trb \"<source path>\" \"<target path>\" --logs \"<logs JSON file path>\"")
	fmt.Println("For usage guide:")
	fmt.Println("\trb")
}
