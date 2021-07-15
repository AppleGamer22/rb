package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	var newLogsPath string
	var fileCount int
	if logsPath != "" {
		stats, err := os.Stat(logsPath)
		if err != nil {
			log.Fatal(err)
		}
		path, count, _, err := rb.GetFilePathsSinceDate(source, target, stats.ModTime())
		if err != nil {
			log.Fatal(err)
		}
		newLogsPath = path
		fileCount = count
	} else {
		var path, count, _, err2 = rb.GetFilePaths(source, target)
		if err2 != nil {
			log.Fatal(err2)
		}
		newLogsPath = path
		fileCount = count
	}
	fmt.Println(newLogsPath, source, target, fileCount, now)
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
